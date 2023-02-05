package database_test

import (
	"testing"

	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var db *database.MongoDb

func TestGoscrapy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database test Suite")
}

var _ = BeforeSuite(func() {

	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	})

	// Getting all config needed for connections and pages login
	setup := config.Get("../data.env", log)

	var err error
	// Starting DB connection with a test database
	db, err = database.Connect(setup.Dburi, setup.Dbuser,
		setup.Dbpass, "mayorista2", log)
	if err != nil {
		log.Panicf("Error database connection %s", err.Error())
	}
	log.Warn("Established database connection")
})

var _ = AfterSuite(func() {
	database.Disconnect()
})
