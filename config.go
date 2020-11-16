package main

import (
	"log"
	"os"

	"github.com/jackgris/goscrapy/wholesalers"
	"github.com/joho/godotenv"
)

var (
	login      string
	user       string
	pass       string
	searchpage string
	dburi      string
	dbuser     string
	dbpass     string
)

func GetConfig() {

	// Getting config data
	err := godotenv.Load("data.env")
	if err != nil {
		log.Fatal("Error loading config file")
	}

	login = os.Getenv("LOGIN")
	user = os.Getenv("USER")
	pass = os.Getenv("PASS")
	searchpage = os.Getenv("SEARCHPAGE")
	dburi = os.Getenv("DBURI")
	dbuser = os.Getenv("DBUSER")
	dbpass = os.Getenv("DBPASS")
}

// This function will get data from the server
func GetWholesalersData() wholesalers.Wholesalers {

	w := wholesalers.Wholesalers{}
	w.Login = login
	w.User = user
	w.Pass = pass
	w.Searchpage = searchpage

	if w.Login == "" || w.Searchpage == "" {
		log.Fatal("Can't get search URI for wholersaler")
	}

	return w
}
