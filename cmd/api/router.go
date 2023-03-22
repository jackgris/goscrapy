package main

import "github.com/gofiber/fiber/v2"

func routes(api fiber.Router) {
	api.Get("/products/:id", GetProductById)
	api.Get("/products", GetAllProducts)
	api.Get("/scraping", Scraper)
	api.Post("/wholesaler", SaveWholesaler)
	api.Get("/wholesaler", GetWholesalers)
	api.Get("/wholesaler/:id", GetWholesaler)
	api.Get("/", Home)
}
