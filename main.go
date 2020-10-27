package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"encoding/json"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/joho/godotenv"

	badger "github.com/dgraph-io/badger/v2"
	"errors"
)

type Product struct {
	MainEntity MainEntity `json:"mainEntityOfPage"`
	Name string `json:"name"`
	Image string `json:"image"`
	Description string `json:"description"`
	Offers Offers
}

type MainEntity struct {
	Id string `json:"@id"`
}

type Offers struct {
	Price string `json:"price"`
	Availability string `json:"availability"`
	InventoryLevel InventoryLevel
}

type InventoryLevel struct {
	Stock string `json:"value"`
}

func main(){

	// Getting config data
	err := godotenv.Load("data.env")

	if err != nil {
		log.Fatal("Error loading config file")
	}

	login := os.Getenv("LOGIN")
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	searchpage := os.Getenv("SEARCHPAGE")


	// Open the Badger database located in the . directory.
	db, err := badger.Open(badger.DefaultOptions("."))
	if err != nil {
		log.Fatal("Error database", err)
	}

	defer db.Close()

	// Starting data collector
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{Delay: 5 * time.Second})

	// With this we know when is the last page of catalog
	end := false

	// Authenticate
	err = c.Post(login, map[string]string{"username": user, "password": pass})
	if err != nil {
		log.Fatal(err)
	}

	// Attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("Response received", r.StatusCode)
	})
	
	c.OnHTML("html", func (e *colly.HTMLElement){
		
		// Goquery selection of the HTMLElement is in e.DOM
		goquerySelection := e.DOM
		
		// Check here, there are or there are no products
		goquerySelection.Find("div.text-center").Each(func (i int, el *goquery.Selection){

			if end {
				return
			}
			not := "No tenemos"
			end = strings.Contains(el.Text(), not)
		})

		if end { return }
		
		// Finding JSON data from scripts
		goquerySelection.Find("script").Each(func (i int, el *goquery.Selection){

			// Create json struct
			p := Product{}
			json.Unmarshal([]byte(el.Text()), &p)
			
			// Is data is ok, saving product on database
			if p.MainEntity.Id != "" {	
				saveData(db, p)
			}
		})
	})
	
	c.OnRequest(func(r *colly.Request){
		log.Println("OnRequest")
	})
	
	// Start scraping FIXME change numbers, this's only for tests
	for i := 63; i < 1000; i++ {
		// Check when there are no products
		if end {
			fmt.Println("Searching end")
			break
		}
		
		num := strconv.Itoa(i)
		URL := searchpage + num

		c.Visit(URL)
	}
	
	fmt.Println(c.String())
}

// Save a product on database
func saveData (db *badger.DB, p Product ){
	
	data, err := json.Marshal(p)
	if err != nil {
		log.Fatal("saveData json marshal ", err)
	}
	
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(p.MainEntity.Id), data)
		return err
	})


	if err != nil {
		log.Fatal("Can't save data: ", err)
	}
}

// Take all products from database
func GetData(db *badger.DB) ([]Product, error) {

	products := []Product{}

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {

			item := it.Item()
			copy := []byte{}
			
			copy, err := item.ValueCopy(copy)
			if err != nil {
				log.Fatal("Error value copy: ", err)
				return err
			}

			product := Product{}

			err = json.Unmarshal(copy, &product)

			if err != nil {
				log.Fatal("JSON unmarshall: ", err)
			}
			
			products = append(products, product)
		}
		return nil
	})

	if err != nil {
		err = errors.New("getData: " + err.Error())
	}
	
	return products, err
}

// Take only one product by ID
func GetDatabyId(db *badger.DB, id string)(Product, error){

	product := Product{}
	// Searching a product
	err := db.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		copy := []byte{}
		copy, err = item.ValueCopy(copy)
		if err != nil {
			return err
		}
		
		err = json.Unmarshal(copy, &product)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.New("getDatabyId: " + err.Error())
	}
	
	return product, err
}
