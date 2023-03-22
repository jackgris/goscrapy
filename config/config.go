package config

import (
	"os"

	"github.com/jackgris/goscrapy/database"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Data struct {
	login        string
	User         string
	Pass         string
	Searchpage   string
	Dburi        string
	Dbuser       string
	Dbpass       string
	NameSaler    string
	Endphrase    string
	Endphrasediv string
}

// Will get all data needed for login on database and web pages
func Get(path string, log *logrus.Logger) Data {

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
	config.Endphrase = os.Getenv("ENDPHRASE")
	config.Endphrasediv = os.Getenv("ENDPHRASEDIV")

	return config
}

// This function will return a wholesaler with all data needed for login, and web navigation
// We need this for extracting and saving data from his web server
func GetWholesalersData(config Data, log *logrus.Logger) database.Wholesalers {

	// FIXME need a more elegant and better solution for this
	w := database.Wholesalers{}
	w.Login = config.login
	w.User = config.User
	w.Pass = config.Pass
	w.Searchpage = config.Searchpage
	w.EndPhrase = config.Endphrase
	w.EndPhraseDiv = config.Endphrasediv
	w.Name = config.NameSaler

	if w.Login == "" || w.Searchpage == "" {
		log.Fatal("Can't get search URI for wholesaler")
	}

	return w
}
