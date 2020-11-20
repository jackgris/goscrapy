package main

import (
	"context"

	"github.com/jackgris/goscrapy/database"
	"github.com/jackgris/goscrapy/wholesalers"
)

func main() {

	GetConfig()
	db, err := database.Connect(dburi, dbuser, dbpass, context.Background())
	if err != nil {
		panic("Error conection: " + err.Error())
	}
	defer database.Disconnect()
	wholesalers.GetData(
		db,
		GetWholesalersData())
}
