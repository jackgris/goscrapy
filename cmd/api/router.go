package main

import "github.com/gofiber/fiber/v2"

func routes(api fiber.Router) {
	api.Get("/products/:id", GetProductById)
	api.Get("/products", GetAllProducts)
	api.Get("/scraping", Scraper)
	api.Get("/", Home)
}
