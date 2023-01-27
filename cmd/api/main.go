package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	// Getting all config needed for connections and pages login
	setup := config.Get("../../data.env")

	// Starting DB connection
	db, err := database.Connect(setup.Dburi, setup.Dbuser, setup.Dbpass, "mayorista")
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()

	app := fiber.New()

	app.Get("/products/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		oId, _ := primitive.ObjectIDFromHex(id)
		product := database.Product{Id_: oId}
		product = db.ReadByMongoId(product)
		empty, _ := primitive.ObjectIDFromHex("000000000000000000000000")
		if product.Id_ == empty {
			r := struct{ Message string }{Message: "BAD ID"}
			return c.JSON(r)
		}
		return c.JSON(product)
	})

	app.Get("/products", func(c *fiber.Ctx) error {
		products := db.ReadByWholesalers(setup.NameSaler)
		return c.JSON(products)
	})

	app.Get("/scraping", func(c *fiber.Ctx) error {
		// Getting and saving data
		err := data.GetData(db, config.GetWholesalersData(setup))
		return err
	})

	log.Fatal(app.Listen(":3000"))
}
