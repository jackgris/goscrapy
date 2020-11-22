package wholesalers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/jackgris/goscrapy/database"

	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// This struct will be use for match JSON data
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

/* With this function we will get data from the web of wholesalers, and save that
information on the database */
func GetData(db database.Database, w Wholesalers) {

	// Starting data collector
	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{Delay: 5 * time.Second})
	if err != nil {
		fmt.Println("Getdata: ", err)
	}
	// With this we know when is the last page of catalog
	end := false

	// Authenticate
	err = c.Post(w.Login, map[string]string{"username": w.User, "password": w.Pass})
	if err != nil {
		log.Fatal("Get Data authenticate: ", err)
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
			_ = json.Unmarshal([]byte(el.Text()), &p)

			// Is data is ok, saving product on database
			if p.MainEntity.Id != "" {
				// saving data here
				//saveData(collection, ctx, p)
				product := database.Product{
					Id:          p.MainEntity.Id,
					Name:        p.Name,
					Image:       p.Image,
					Description: p.Description,
					Price:       p.Offers.Price,
					Stock:       p.Offers.InventoryLevel.Stock,
					Wholesaler:  "acabajo", // FIXME this data, will not be in code
				}
				err := db.Create(product)
				if err != nil {
					fmt.Println("Can´t save product: ", err)
				}
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("OnRequest")
	})

	// FIXME Start scraping change numbers, this's only for tests
	for i := 63; i < 1000; i++ {
		// Check when there are no products
		if end {
			fmt.Println("Searching end")
			break
		}

		num := strconv.Itoa(i)
		URL := w.Searchpage + num

		err := c.Visit(URL)
		if err != nil {
			fmt.Println("Error visiting site: ", err)
		}
	}

	fmt.Println(c.String())
}
