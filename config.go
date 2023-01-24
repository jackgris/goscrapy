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

// Will get all data needed for login on database and web pages
func getConfig() {

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

// This function will return a wholesaler with all data needed for login, and web navigation
// We need this for extracting and saving data from his web server
func GetWholesalersData() wholesalers.Wholesalers {

	// FIXME need a more elegant and better solution for this
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
