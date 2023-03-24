package main

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Home(c *fiber.Ctx) error {
	database.Db.Log.Info("Until now, this is only for test propuse. Need create an html template.")
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
	// Get all the wholesalers save in database.
	wholesalers := database.Db.FindWholesalers()

	// If don't have any wholesaler saved in our database, we can try use data from a config file.
	if len(wholesalers) == 0 {
		wholesalers = append(wholesalers, config.GetWholesalersData(setup, database.Db.Log))
	}

	var wg sync.WaitGroup
	wg.Add(len(wholesalers))

	// Getting and saving all product data from all wholesaler running in parallel.
	for _, wholesaler := range wholesalers {
		go func(w database.Wholesalers) {
			defer wg.Done()
			err := data.GetData(database.Db, w, database.Db.Log)
			if err != nil {
				database.Db.Log.Fatal(err)
			}
		}(wholesaler)
	}
	wg.Wait()

	// Finish scrape all wholesaler without any problem.
	database.Db.Log.Warn("Scrape finished")

	return c.Status(http.StatusOK).SendString("All done")
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

func GetWholesaler(c *fiber.Ctx) error {

	var wholesaler database.Wholesalers
	id := c.Params("id")
	oId, _ := primitive.ObjectIDFromHex(id)
	w := database.Wholesalers{Id: oId}
	wholesaler = database.Db.GetWhosalerById(w)
	empty, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	if wholesaler.Id == empty {
		r := struct{ Message string }{Message: "BAD ID"}
		return c.JSON(r)
	}
	return c.JSON(wholesaler)
}

func GetWholesalers(c *fiber.Ctx) error {
	wholesalers := database.Db.FindWholesalers()
	return c.JSON(wholesalers)
}

func UpdateWholesaler(c *fiber.Ctx) error {

	var ws database.Wholesalers
	if err := c.BodyParser(&ws); err != nil {
		database.Db.Log.Warn("UpdateWholesaler:", err)
		return c.Status(http.StatusBadRequest).SendString("Bad request")
	}

	if ok := checkWSaler(ws); !ok {
		database.Db.Log.Warn("UpdateWholesaler:", err)
		return c.Status(http.StatusNotAcceptable).SendString("Incomplete data")
	}
	if err := database.Db.UpdateWholesaler(ws); err != nil {
		database.Db.Log.Warn("UpdateWholesaler:", err)
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
