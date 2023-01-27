package database_test

import (
	"testing"

	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var db *database.MongoDb

func TestGoscrapy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Goscrapy Suite")
}

var _ = BeforeSuite(func() {

	// Getting all config needed for connections and pages login
	setup := config.Get("../data.env")

	var err error
	// Starting DB connection
	db, err = database.Connect(setup.Dburi, setup.Dbuser, setup.Dbpass, "mayorista2")
	if err != nil {
		panic("Error database connection: " + err.Error())
	}
})

var _ = AfterSuite(func() {
	database.Disconnect()
})
