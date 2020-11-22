package main

import (
	"github.com/jackgris/goscrapy/database"
	"github.com/jackgris/goscrapy/wholesalers"
)

func main() {

	// Getting all config needed for connections and pages login
	GetConfig()
	// Starting DB connection
	db, err := database.Connect(dburi, dbuser, dbpass)
	if err != nil {
		panic("Error conection: " + err.Error())
	}
	defer database.Disconnect()
	// Getting and saving data
	wholesalers.GetData(
		db,
		GetWholesalersData())
}
