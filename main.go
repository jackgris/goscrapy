package main

import (
	"github.com/jackgris/goscrapy/database"
	"github.com/jackgris/goscrapy/wholesalers"
)

func main() {

	// Getting all config needed for connections and pages login
	getConfig()
	// Starting DB connection
	db, err := database.Connect(dburi, dbuser, dbpass)
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
	defer database.Disconnect()
	// Getting and saving data
	wholesalers.GetData(
		db,
		GetWholesalersData())
}
