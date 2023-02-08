package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Home(c *fiber.Ctx) error {
	database.Db.Log.Info("Until now, this is only for test propuse")
	r := struct{ Message string }{Message: "THIS HOME"}
	return c.JSON(r)
}

func GetProductById(c *fiber.Ctx) error {
	id := c.Params("id")
	oId, _ := primitive.ObjectIDFromHex(id)
	product := database.Product{Id_: oId}
	product = database.Db.ReadByMongoId(product)
	empty, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	if product.Id_ == empty {
		r := struct{ Message string }{Message: "BAD ID"}
		return c.JSON(r)
	}
	return c.JSON(product)
}

func GetAllProducts(c *fiber.Ctx) error {
	products := database.Db.ReadByWholesalers(setup.NameSaler)
	return c.JSON(products)
}

func Scraper(c *fiber.Ctx) error {
	// Getting and saving data
	err := data.GetData(database.Db, config.GetWholesalersData(setup, database.Db.Log), database.Db.Log)
	return err
}
