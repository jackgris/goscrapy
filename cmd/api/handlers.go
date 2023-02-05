package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Home(c *fiber.Ctx) error {
	db.Log.Info("Until now, this is only for test propuse")
	r := struct{ Message string }{Message: "THIS HOME"}
	return c.JSON(r)
}

func getProductById(c *fiber.Ctx) error {
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
}

func getAllProducts(c *fiber.Ctx) error {
	products := db.ReadByWholesalers(setup.NameSaler)
	return c.JSON(products)
}

func scraper(c *fiber.Ctx) error {
	// Getting and saving data
	err := data.GetData(db, config.GetWholesalersData(setup, db.Log), db.Log)
	return err
}
