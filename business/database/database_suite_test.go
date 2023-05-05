package database_test

import (
	"testing"

	"github.com/jackgris/goscrapy/business/database"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/foundation/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoscrapy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database test Suite")
}

var _ = BeforeSuite(func() {

	log := logger.New()

	// Getting all config needed for connections and pages login
	setup := config.Get("../data.env", log)

	var err error
	// Starting DB connection with a test database
	_, err = database.Connect(setup.Dburi, setup.Dbuser,
		setup.Dbpass, "mayorista2", log)
	if err != nil {
		log.Panicf("Error database connection %s", err.Error())
	}
	log.Warn("Established database connection")
})

var _ = AfterSuite(func() {
	database.Disconnect()
})
