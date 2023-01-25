package main

import (
	"github.com/jackgris/goscrapy/data"
	"github.com/jackgris/goscrapy/database"
)

func main() {

	// Getting all config needed for connections and pages login
	config := getConfig()
	// Starting DB connection
	db, err := database.Connect(config.dburi, config.dbuser, config.dbpass)
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()
	// Getting and saving data
	data.GetData(db, getWholesalersData(config))
}
