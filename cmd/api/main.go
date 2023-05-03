package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	v1 "github.com/jackgris/goscrapy/cmd/api/handlers/v1"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/database"
	"github.com/sirupsen/logrus"
)

var err error
var setup config.Data

func main() {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	})

	// Getting all config needed for connections and pages login
	setup = config.Get("data.env", log)

	// Starting DB connection
	db, err := database.Connect(setup.Dburi, setup.Dbuser,
		setup.Dbpass, "mayorista", log)

	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()

	app := fiber.New()

	// Initialize default config
	conf := logger.ConfigDefault
	conf.Format = "${time} | ${status} | ${latency} | IP ${ip}  | Method  ${method}| Route ${path}\nResponse: ${resBody}\n"
	// 13:54:41 | 200 |    81ms |       127.0.0.1 | GET     | /spreadsheet
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

	v1.Routes(app, cfg)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

	// Waiting for start shutting down
	serverShutdown.Wait()
	log.Warn("Running cleanup tasks...")
}
