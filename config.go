package main

import (
	"log"
	"os"

	"github.com/jackgris/goscrapy/data"
	"github.com/joho/godotenv"
)

type configData struct {
	login      string
	user       string
	pass       string
	searchpage string
	dburi      string
	dbuser     string
	dbpass     string
}

// Will get all data needed for login on database and web pages
func getConfig() configData {

	// Getting config data
	err := godotenv.Load("data.env")
	if err != nil {
		log.Fatal("Error loading config file")
	}
	config := configData{}
	config.login = os.Getenv("LOGIN")
	config.user = os.Getenv("USER")
	config.pass = os.Getenv("PASS")
	config.searchpage = os.Getenv("SEARCHPAGE")
	config.dburi = os.Getenv("DBURI")
	config.dbuser = os.Getenv("DBUSER")
	config.dbpass = os.Getenv("DBPASS")

	return config
}

// This function will return a wholesaler with all data needed for login, and web navigation
// We need this for extracting and saving data from his web server
func getWholesalersData(config configData) data.Wholesalers {

	// FIXME need a more elegant and better solution for this
	w := data.Wholesalers{}
	w.Login = config.login
	w.User = config.user
	w.Pass = config.pass
	w.Searchpage = config.searchpage

	if w.Login == "" || w.Searchpage == "" {
		log.Fatal("Can't get search URI for wholesaler")
	}

	return w
}
