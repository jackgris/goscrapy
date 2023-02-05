package database

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Interface useful for mocking test, or if need change database
type Database interface {
	Create(p Product) error
	ReadById(p Product) Product
	ReadByWholesalers(name string) []Product
	Delete(p Product) error
}

// My database struct
type MongoDb struct {
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
	name   string
	Log    *logrus.Logger
}

// vars needed for only created one instance for my access to the database
var (
	once sync.Once
	db   *MongoDb
)

// Connecting to the database, only one instance will be create, to connect with mongo database, we need
// the URI where are the DB, the user name and password, and will return the instance with an active connection
// or an error
func Connect(dburi, dbuser, dbpass, name string, log *logrus.Logger) (*MongoDb, error) {

	var err error
	once.Do(func() {
		db = new(MongoDb)
		db.Log = log
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		db.cancel = cancel
		db.ctx = ctx
		option := options.Client().ApplyURI(dburi)
		credentials := options.Credential{Username: dbuser, Password: dbpass}
		option.SetAuth(credentials)

		db.client, err = mongo.Connect(db.ctx, option)
		if err != nil {
			log.Panicf("Can't connect with db: %s", err.Error())
		}
		err = db.client.Ping(db.ctx, readpref.Primary())

	})

	db.name = name

	return db, err
}

// Closing the database connection, correctly
func Disconnect() {
	db.Log.Warn("Disconnecting from database")
	defer db.cancel()
	if err := db.client.Disconnect(db.ctx); err != nil {
		db.Log.Panicf("Error trying close connection db: %s", err.Error())
	}
}
