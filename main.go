package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
)

func main() {

	// Getting all config needed for connections and pages login
	setup := config.Get("data.env")

	// Starting DB connection
	db, err := database.Connect(setup.Dburi, setup.Dbuser, setup.Dbpass, "mayorista")
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()

	app := fiber.New()

	app.Get("/products", func(c *fiber.Ctx) error {
		return nil
	})

	app.Get("/scraping", func(c *fiber.Ctx) error {
		// Getting and saving data
		err := data.GetData(db, config.GetWholesalersData(setup))
		return err
	})

	log.Fatal(app.Listen(":3000"))
}
