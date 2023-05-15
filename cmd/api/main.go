package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackgris/goscrapy/business/database"
	v1 "github.com/jackgris/goscrapy/cmd/api/handlers/v1"
	"github.com/jackgris/goscrapy/config"
	logs "github.com/jackgris/goscrapy/foundation/logger"
)

func main() {
	log := logs.New()

	// Getting all config needed for connections and pages login
	setup := config.Get("data.env", log)

	credentials := database.Credentials{User: setup.Dbuser, Password: setup.Dbpass}
	// Starting DB connection
	db, err := database.Connect(setup.Dburi, "mayorista", log, credentials)

	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()

	app := fiber.New()

	// Initialize default config
	conf := logger.ConfigDefault
	conf.Format = "${time} | ${status} | ${latency} | IP ${ip}  | Method  ${method}| Route ${path} Error: ${error}\n"
	app.Use(logger.New(conf))

	// Will wait for signal interrupt, to wait for a while and clean all the pending tasks.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var serverShutdown sync.WaitGroup

	go func() {
		<-c
		log.Warn("Gracefully shutting down...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(60 * time.Second)
	}()

	cfg := v1.Config{
		Log:   log,
		Db:    db,
		Setup: &setup,
	}

	// Add all the handler for API v1
	v1.Routes(app, cfg)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

	// Waiting for start shutting down
	serverShutdown.Wait()
	log.Warn("Running cleanup tasks...")
}
