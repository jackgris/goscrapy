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

	
	err := godotenv.Load("data.env")

	if err != nil {
		log.Fatal("Error loading config file")
	}

	login := os.Getenv("LOGIN")
	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	searchpage := os.Getenv("SEARCHPAGE")
	
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
		
		// I need find some better filter than count := 0 = 11	:D
		goquerySelection.Find("script").Each(func (i int, el *goquery.Selection){

			// Create json struct
			var p Product
			json.Unmarshal([]byte(el.Text()), &p)
			json, err := json.MarshalIndent(p, "", "\t")
			if err != nil {
				log.Fatal(err)
			}

			// FIXME save in database or file
			if p.MainEntity.Id != "" {
				fmt.Println(string(json))				
			}
		})
	})
	
	c.OnRequest(func(r *colly.Request){
		log.Println("OnRequest")
	})
	
	// Start scraping FIXME change numbers, this's only for tests
	for i := 62; i < 1000; i++ {
		// Check when there are no products
		if end {
			fmt.Println("Searching end")
			return
		}
		
		num := strconv.Itoa(i)
		URL := searchpage + num
		fmt.Println("URL: ", URL)
		c.Visit(URL)
	}

	// End
	fmt.Println(c.String())
}
