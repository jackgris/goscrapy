package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
)

func main() {

	// Getting all config needed for connections and pages login
	config := getConfig()
	// Starting DB connection
	db, err := database.Connect(config.dburi, config.dbuser, config.dbpass)
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()

	app := fiber.New()

	app.Get("/scraping", func(c *fiber.Ctx) error {
		// Getting and saving data
		err := data.GetData(db, getWholesalersData(config))
		return err
	})

	log.Fatal(app.Listen(":3000"))
}
