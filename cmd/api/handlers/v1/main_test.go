package v1_test

import (
	"fmt"
	"testing"

	"github.com/jackgris/goscrapy/business/database"
	"github.com/jackgris/goscrapy/business/dbtest"
	v1 "github.com/jackgris/goscrapy/cmd/api/handlers/v1"
	"github.com/jackgris/goscrapy/config"
	"github.com/jackgris/goscrapy/foundation/docker"
	logs "github.com/jackgris/goscrapy/foundation/logger"
	"github.com/sirupsen/logrus"
)

var c *docker.Container
var cfg *v1.Config
var log *logrus.Logger
var db *database.MongoDb

func TestMain(m *testing.M) {

	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	log = logs.New()

	// Getting all config needed for connections and pages login
	setup := config.Data{
		Dburi:  "mongodb://" + c.Host,
		Dbuser: "admin",
		Dbpass: "admin",
	}

	// Starting DB connection
	db, err = database.Connect(setup.Dburi, "mayorista", log)
	if err != nil {
		panic("Error database connection: " + err.Error())
	}

	cfg = &v1.Config{
		Log:   log,
		Db:    db,
		Setup: &setup,
	}

	m.Run()
}
