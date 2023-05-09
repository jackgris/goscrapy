package v1

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/jackgris/goscrapy/business/data"
	"github.com/jackgris/goscrapy/business/database"
	"github.com/jackgris/goscrapy/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Home(c *fiber.Ctx) error {
	r := struct{ Message string }{Message: "THIS HOME"}
	c.Status(http.StatusOK)
	return c.JSON(r)
}

func GetProductById(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {

		id := c.Params("id")
		oId, _ := primitive.ObjectIDFromHex(id)
		product := database.Product{Id_: oId}
		product = cfg.Db.ReadByMongoId(product)
		empty, _ := primitive.ObjectIDFromHex("000000000000000000000000")
		if product.Id_ == empty {
			r := struct{ Message string }{Message: "BAD ID"}
			return c.JSON(r)
		}
		return c.JSON(product)
	}

	return fn
}

func GetProductsByWholesaler(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		name := c.Params("wholesaler")
		products := cfg.Db.ReadByWholesalers(name)

		return c.JSON(products)
	}
	return fn
}

func GetAllProducts(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		products := cfg.Db.GetAllProducts()

		return c.JSON(products)
	}
	return fn
}

func Scraper(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		// Get all the wholesalers save in database.
		wholesalers := cfg.Db.FindWholesalers()

		// If don't have any wholesaler saved in our database, we can try use data from a config file.
		if len(wholesalers) == 0 {

			wholesalers = append(wholesalers, config.GetWholesalersData(*cfg.Setup, database.Db.Log))
		}

		var wg sync.WaitGroup
		wg.Add(len(wholesalers))

		// Getting and saving all product data from all wholesaler running in parallel.
		for _, wholesaler := range wholesalers {
			go func(w database.Wholesalers) {
				defer wg.Done()
				err := data.GetData(cfg.Db, w, cfg.Log)
				if err != nil {
					cfg.Log.Fatal(err)
				}
			}(wholesaler)
		}
		wg.Wait()

		// Finish scrape all wholesaler without any problem.
		cfg.Log.Warn("Scrape finished")

		return c.Status(http.StatusOK).SendString("All done")
	}
	return fn
}

func SaveWholesaler(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		var ws database.Wholesalers
		if err := c.BodyParser(&ws); err != nil {
			cfg.Log.Warn("SaveWhosaler:", err)
			return c.Status(http.StatusBadRequest).SendString("Bad request")
		}

		if ok := checkWSaler(ws); !ok {
			cfg.Log.Warn("SaveWhosaler: some data are empty")
			return c.Status(http.StatusNotAcceptable).SendString("Incomplete data")
		}
		if err := cfg.Db.InsertWholesaer(ws); err != nil {
			cfg.Log.Warn("SaveWhosaler:", err)
			return c.Status(http.StatusInternalServerError).SendString("Can't saved entity")
		}

		return c.Status(http.StatusOK).SendString("All is right")

	}
	return fn
}

func GetWholesaler(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		var wholesaler database.Wholesalers
		id := c.Params("id")
		oId, _ := primitive.ObjectIDFromHex(id)
		w := database.Wholesalers{Id: oId}
		wholesaler = cfg.Db.GetWhosalerById(w)
		empty, _ := primitive.ObjectIDFromHex("000000000000000000000000")
		if wholesaler.Id == empty {
			r := struct{ Message string }{Message: "BAD ID"}
			return c.JSON(r)
		}
		return c.JSON(wholesaler)
	}
	return fn
}

func GetWholesalers(cfg Config) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		wholesalers := cfg.Db.FindWholesalers()
		return c.JSON(wholesalers)
	}
	return fn
}

func UpdateWholesaler(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {

		var ws database.Wholesalers
		if err := c.BodyParser(&ws); err != nil {
			cfg.Log.Warn("UpdateWholesaler:", err)
			return c.Status(http.StatusBadRequest).SendString("Bad request")
		}

		if ok := checkWSaler(ws); !ok {
			cfg.Log.Warn("UpdateWholesaler: Some necesary data is empty")
			return c.Status(http.StatusNotAcceptable).SendString("Incomplete data")
		}
		if err := cfg.Db.UpdateWholesaler(ws); err != nil {
			cfg.Log.Warn("UpdateWholesaler:", err)
			return c.Status(http.StatusInternalServerError).SendString("Can't saved entity")
		}

		return c.Status(http.StatusOK).SendString("All is right")
	}
	return fn
}

func checkWSaler(ws database.Wholesalers) (ok bool) {
	ok = true
	if ws.Login == "" || ws.Pass == "" || ws.Name == "" || ws.Searchpage == "" || ws.User == "" {
		ok = false
	}
	return ok
}

func ShowSameProducts(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		// Get all products saved in our database.
		products := cfg.Db.GetAllProducts()
		var result [][]database.Product
		for _, p := range products {

			// Clean from the list products from the same owner.
			if p.Wholesaler == cfg.Setup.NameSaler {

				// And only add result with at least one match.
				r := cfg.Db.SearchSimilars(p)
				if len(r) > 1 {
					result = append(result, r)
				}
			}
		}
		return c.JSON(result)
	}
	return fn
}

func ComparePricesSameWholesaler(cfg Config) fiber.Handler {

	fn := func(c *fiber.Ctx) error {
		name := c.Params("wholesaler")
		products := data.ReadCSV("./business/data/csv/"+name+".csv", cfg.Log)
		var result []PriceCompared
		for _, p := range products {
			r := cfg.Db.SearchByName(p)
			if len(r) >= 1 {
				list := []PriceWeb{}
				for _, webP := range r {
					web := PriceWeb{
						Name:  webP.Name,
						Price: webP.Price[len(webP.Price)-1].Price,
						Owner: webP.Wholesaler,
					}
					list = append(list, web)
				}

				pC := PriceCompared{
					Name:  p.Name,
					Price: p.Price[len(p.Price)-1].Price,
					Webs:  list,
				}

				if pC.Price < pC.Webs[len(pC.Webs)-1].Price {
					cfg.Log.Warnf("Need update price of %s actual prices %v and new is %v.",
						pC.Name, pC.Price, pC.Webs[len(pC.Webs)-1].Price)
				}

				result = append(result, pC)
			}
		}

		return c.JSON(result)
	}
	return fn
}

func NeedUpdatePricesSameWholesaler(cfg Config) fiber.Handler {
	fn := func(c *fiber.Ctx) error {

		name := c.Params("wholesaler")
		discount := c.Params("discount", "")
		percentage := float64(1)

		num, err := strconv.ParseFloat(discount, 64)

		if err == nil && num != 0 {
			num = num / 100
			percentage -= num
		}

		products := data.ReadCSV("./business/data/csv/"+name+".csv", cfg.Log)
		var result []PriceCompared
		for _, p := range products {
			r := cfg.Db.SearchByName(p)
			if len(r) >= 1 {
				list := []PriceWeb{}
				for _, webP := range r {
					web := PriceWeb{
						Name:  webP.Name,
						Price: webP.Price[len(webP.Price)-1].Price * percentage,
						Owner: webP.Wholesaler,
					}
					list = append(list, web)
				}

				pC := PriceCompared{
					Name:  p.Name,
					Price: p.Price[len(p.Price)-1].Price,
					Webs:  list,
				}

				if pC.Price < pC.Webs[len(pC.Webs)-1].Price {
					cfg.Log.Warnf("Need update price of %s actual prices %v and new is %v.",
						pC.Name, pC.Price, pC.Webs[len(pC.Webs)-1].Price)

					result = append(result, pC)
				}
			}
		}

		return c.JSON(result)
	}
	return fn
}

func CreateXlsxFile(cfg Config) fiber.Handler {
	fn := func(c *fiber.Ctx) error {
		products := cfg.Db.GetAllProducts()
		path := data.WriteXlsx(".", cfg.Log, products)
		if path != "" {
			c.Attachment("./NewPrices.xlsx")
			return c.SendFile(path)
		}
		cfg.Log.Warn("Can't return xlsx file")
		return nil
	}
	return fn
}
