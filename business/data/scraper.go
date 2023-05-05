package data

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/jackgris/goscrapy/business/database"
	"github.com/sirupsen/logrus"

	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type SaveUser interface {
	Create(p database.Product) error
}

// GetData will get data from the web of wholesalers, and save that information on the database
func GetData(db SaveUser, w database.Wholesalers, log *logrus.Logger) error {

	// Starting data collector
	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{
		Delay:        5 * time.Second,
		DomainRegexp: w.Searchpage + "*",
	})
	if err != nil {
		log.Error("Getdata: " + err.Error())
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
		log.Println("Response received: ", r.StatusCode, " URL: ", r.Request.URL)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {

		// Goquery selection of the HTMLElement is in e.DOM
		goquerySelection := e.DOM

		// Check here, there are or there are no products
		goquerySelection.Find(w.EndPhraseDiv).Each(func(i int, el *goquery.Selection) {

			if end {
				return
			}
			not := w.EndPhrase
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
				price, err := strconv.ParseFloat(p.Offers.Price, 64)
				if err != nil {
					log.Warn("Error parsing price: ", err)
				} else {
					var value database.Value
					value.Date = time.Now()
					value.Price = price
					values := []database.Value{value}

					product := database.Product{
						Id:          p.MainEntity.Id,
						Name:        p.Name,
						Image:       p.Image,
						Description: p.Description,
						Price:       values,
						Stock:       p.Offers.InventoryLevel.Stock,
						Wholesaler:  w.Name,
					}
					err = db.Create(product)
					if err != nil {
						log.Warn("Can´t save product: ", err)
					}
				}
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {})

	// Start scraping change numbers, this's only for tests
	for i := 1; i < 1000; i++ {
		// Check when there are no products
		if end {
			log.Debug("Searching end")
			break
		}

		num := strconv.Itoa(i)
		URL := w.Searchpage + num

		err := c.Visit(URL)
		if err != nil {
			log.Debug("Error visiting site: ", err)
		}
	}

	log.Debug(c.String())

	return nil
}
