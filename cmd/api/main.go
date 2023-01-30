package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/database"
)

var db *database.MongoDb
var err error
var setup config.Data

func main() {

	// Getting all config needed for connections and pages login
	setup = config.Get("../../data.env")

	// Starting DB connection
	db, err = database.Connect(setup.Dburi, setup.Dbuser, setup.Dbpass, "mayorista")
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()

	app := fiber.New()
	app.Get("/products/:id", getProductById)
	app.Get("/products", getAllProducts)
	app.Get("/scraping", scraper)
	log.Fatal(app.Listen(":3000"))
}
