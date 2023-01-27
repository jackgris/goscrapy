package config

import (
	"log"
	"os"

	"github.com/jackgris/goscrapy/data"
	"github.com/joho/godotenv"
)

type Data struct {
	login      string
	User       string
	Pass       string
	Searchpage string
	Dburi      string
	Dbuser     string
	Dbpass     string
	NameSaler  string
}

// Will get all data needed for login on database and web pages
func Get(path string) Data {

	// Getting config data
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading config file")
	}
	config := Data{}
	config.login = os.Getenv("LOGIN")
	config.User = os.Getenv("USER")
	config.Pass = os.Getenv("PASS")
	config.Searchpage = os.Getenv("SEARCHPAGE")
	config.Dburi = os.Getenv("DBURI")
	config.Dbuser = os.Getenv("DBUSER")
	config.Dbpass = os.Getenv("DBPASS")
	config.NameSaler = os.Getenv("WHOLESALER")

	return config
}

// This function will return a wholesaler with all data needed for login, and web navigation
// We need this for extracting and saving data from his web server
func GetWholesalersData(config Data) data.Wholesalers {

	// FIXME need a more elegant and better solution for this
	w := data.Wholesalers{}
	w.Login = config.login
	w.User = config.User
	w.Pass = config.Pass
	w.Searchpage = config.Searchpage

	if w.Login == "" || w.Searchpage == "" {
		log.Fatal("Can't get search URI for wholesaler")
	}

	return w
}
