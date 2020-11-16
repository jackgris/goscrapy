package wholesalers

import (
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"

	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strconv"
	"strings"
	"time"
)

type Product struct {
	MainEntity  MainEntity `json:"mainEntityOfPage"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Offers      Offers
}

type MainEntity struct {
	Id string `json:"@id"`
}

type Offers struct {
	Price          string `json:"price"`
	Availability   string `json:"availability"`
	InventoryLevel InventoryLevel
}

type InventoryLevel struct {
	Stock string `json:"value"`
}

type Wholesalers struct {
	Login      string
	User       string
	Pass       string
	Searchpage string
}

func GetData(dburi, dbuser, dbpass string, w Wholesalers) {

	// Get database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	option := options.Client().ApplyURI(dburi)
	credentials := options.Credential{Username: dbuser, Password: dbpass}
	option.SetAuth(credentials)

	client, err := mongo.Connect(ctx, option)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Can't find database: " + err.Error())
	}

	collection := client.Database("mayorista").Collection("productos")

	// Starting data collector
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Delay: 5 * time.Second})

	// With this we know when is the last page of catalog
	end := false

	// Authenticate
	err = c.Post(w.Login, map[string]string{"username": w.User, "password": w.Pass})
	if err != nil {
		log.Fatal(err)
	}

	// Attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("Response received", r.StatusCode)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {

		// Goquery selection of the HTMLElement is in e.DOM
		goquerySelection := e.DOM

		// Check here, there are or there are no products
		goquerySelection.Find("div.text-center").Each(func(i int, el *goquery.Selection) {

			if end {
				return
			}
			not := "No tenemos"
			end = strings.Contains(el.Text(), not)
		})

		if end {
			return
		}

		// Finding JSON data from scripts
		goquerySelection.Find("script").Each(func(i int, el *goquery.Selection) {

			// Create json struct
			p := Product{}
			json.Unmarshal([]byte(el.Text()), &p)

			// Is data is ok, saving product on database
			if p.MainEntity.Id != "" {
				// saving data here
				saveData(collection, ctx, p)
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("OnRequest")
	})

	// FIXME Start scraping FIXME change numbers, this's only for tests
	for i := 63; i < 1000; i++ {
		// Check when there are no products
		if end {
			fmt.Println("Searching end")
			break
		}

		num := strconv.Itoa(i)
		URL := w.Searchpage + num

		c.Visit(URL)
	}

	fmt.Println(c.String())
}

// Save a product on database
func saveData(collection *mongo.Collection, ctx context.Context, p Product) {

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		log.Fatal("Error inserting json: ", err.Error())
	}

}

// FIXME Take all products from database
func GetAll(db string) ([]Product, error) {
	log.Fatal("Not Implemented getdata wholersaler")
	return []Product{}, nil
}

// FIXME NOT IMPLEMENTED Take only one product by ID
func GetById(db, id string) (Product, error) {
	log.Fatal("Not implemented getdatabyId wholersaler")
	return Product{}, nil
}
