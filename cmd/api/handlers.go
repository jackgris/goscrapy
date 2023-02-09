package main

import (
	"net/http"

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

func SaveWholesaler(c *fiber.Ctx) error {

	var ws database.Wholesalers
	if err := c.BodyParser(&ws); err != nil {
		database.Db.Log.Warn("SaveWhosaler:", err)
		return c.Status(http.StatusBadRequest).SendString("Bad request")
	}

	if ok := checkWSaler(ws); !ok {
		database.Db.Log.Warn("SaveWhosaler:", err)
		return c.Status(http.StatusNotAcceptable).SendString("Incomplete data")
	}
	if err := database.Db.InsertWholesaer(ws); err != nil {
		database.Db.Log.Warn("SaveWhosaler:", err)
		return c.Status(http.StatusInternalServerError).SendString("Can't saved entity")
	}

	return c.Status(http.StatusOK).SendString("All is right")
}

func checkWSaler(ws database.Wholesalers) (ok bool) {
	ok = true
	if ws.Login == "" || ws.Pass == "" || ws.Name == "" || ws.Searchpage == "" || ws.User == "" {
		ok = false
	}
	return ok
}
