package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/business/database"
	"github.com/jackgris/goscrapy/config"
	"github.com/sirupsen/logrus"
)

type PriceCompared struct {
	Name  string
	Price float64
	Webs  []PriceWeb
}

type PriceWeb struct {
	Name  string
	Price float64
	Owner string
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log   *logrus.Logger
	Db    *database.MongoDb
	Setup *config.Data
}

// Routes binds all the version 1 routes.
func Routes(app *fiber.App, cfg Config) func(routes fiber.Router) {

	const version = "v1"

	routes := routes(version, cfg)
	app.Route("/", routes, "v1")

	return routes
}

// routes return the function to setup routes in Fiber
func routes(version string, cfg Config) func(fiber.Router) {

	fn := func(routes fiber.Router) {
		routes.Get(version+"/products/:id", GetProductById(cfg))
		routes.Get(version+"/products/wholesaler/:wholesaler", GetProductsByWholesaler(cfg))
		routes.Get(version+"/products", GetAllProducts(cfg))
		routes.Get(version+"/scraping", Scraper(cfg))
		routes.Get(version+"/similars", ShowSameProducts(cfg))
		routes.Get(version+"/compare/:wholesaler", ComparePricesSameWholesaler(cfg))
		routes.Get(version+"/needupdate/:wholesaler/:discount", NeedUpdatePricesSameWholesaler(cfg))
		routes.Get(version+"/spreadsheet", CreateXlsxFile(cfg))
		routes.Post(version+"/wholesaler", SaveWholesaler(cfg))
		routes.Get(version+"/wholesaler", GetWholesalers(cfg))
		routes.Get(version+"/wholesaler/:id", GetWholesaler(cfg))
		routes.Post(version+"/wholesaler/:id", UpdateWholesaler(cfg))
		routes.Get(version+"/", Home)
	}
	return fn
}
