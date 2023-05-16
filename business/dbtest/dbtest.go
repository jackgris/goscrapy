package dbtest

import (
	"fmt"
	"testing"

	"github.com/jackgris/goscrapy/business/database"
	"github.com/jackgris/goscrapy/foundation/docker"
	"github.com/jackgris/goscrapy/foundation/logger"
	"github.com/sirupsen/logrus"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// StartDB starts a database instance.
func StartDB() (*docker.Container, error) {

	image := "mongo:6.0.5"
	port := "27017"
	args := []string{
		// This are examples of flags you can add to run the database
		// "--name", "goscrapy-mongodb-test",
		// "-e", "MONGO_INITDB_DATABASE=admin",
		// "-e", "MONGODB_INITDB_ROOT_USERNAME=admin",
		// "-e", "MONGODB_INITDB_ROOT_PASSWORD=admin",
	}

	c, err := docker.StartContainer(image, port, args...)
	if err != nil {
		return nil, fmt.Errorf("starting container: %w", err)
	}

	fmt.Printf("Image:       %s\n", image)
	fmt.Printf("ContainerID: %s\n", c.ID)
	fmt.Printf("Host:        %s\n", c.Host)

	return c, nil
}

// StopDB stops a running database instance.
func StopDB(c *docker.Container) {
	_ = docker.StopContainer(c.ID)
	fmt.Println("Stopped:", c.ID)
}

// NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewUnit(t *testing.T, c *docker.Container) (*logrus.Logger, *database.MongoDb, func()) {

	log := logger.New()

	dburi := "mongodb://" + c.Host
	dbname := "mayorista"

	// Starting DB connection
	db, err := database.Connect(dburi, dbname, log)

	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Log("Waiting for database to be ready ...")

	// Use the MongoClient instance to connect to the Mongo database.
	err = db.Ping()
	if err != nil {
		t.Fatalf("Can't connect to the database: %s", err)
	}

	t.Log("Database ready")
	t.Log("Migrate and seed database ...")

	// Migrate data from a previous backup
	err = docker.Restore(c.ID, dbname)
	if err != nil {
		t.Fatalf("Error while restoring backup: %s", err)
	}

	t.Log("Ready for testing ...")

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		database.Disconnect()
	}

	return log, db, teardown
}

// Test owns state for running and shutting down tests.
type Test struct {
	DB       *database.MongoDb
	Log      *logrus.Logger
	Teardown func()

	t *testing.T
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T, c *docker.Container) *Test {
	log, db, teardown := NewUnit(t, c)

	test := Test{
		DB:       db,
		Log:      log,
		t:        t,
		Teardown: teardown,
	}

	return &test
}
