package database

import (
	"context"
	"sync"

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
	Db   *MongoDb
)

// Ping sends a ping command to verify that the client can connect to the deployment.
func (db *MongoDb) Ping() error {
	return db.client.Ping(db.ctx, nil)
}

type Credentials struct {
	User     string
	Password string
	DbName   string
}

func loadCredentials(credentials ...Credentials) Credentials {
	credential := Credentials{}

	for _, c := range credentials {
		credential.User = c.User
		credential.Password = c.Password
		credential.DbName = c.DbName
	}
	return credential
}

// Connect give you a connection to the database, only one instance will be create, to connect with mongo database, we need
// the URI where are the DB, the user name and password, and will return the instance with an active connection
// or an error
func Connect(dburi, name string, log *logrus.Logger, credentials ...Credentials) (*MongoDb, error) {

	cred := loadCredentials(credentials...)
	var err error
	once.Do(func() {
		Db = new(MongoDb)
		Db.Log = log
		ctx, cancel := context.WithCancel(context.Background())
		Db.cancel = cancel
		Db.ctx = ctx
		option := options.Client().ApplyURI(dburi)

		if cred.User != "" {
			credentials := options.Credential{Username: cred.User, Password: cred.Password, PasswordSet: true}
			option.SetAuth(credentials)
		}

		Db.client, err = mongo.Connect(Db.ctx, option)
		if err != nil {
			log.Panicf("Can't connect with db: %s", err.Error())
		}
		err = Db.client.Ping(Db.ctx, readpref.Primary())

	})

	Db.name = name

	return Db, err
}

// Closing the database connection, correctly
func Disconnect() {
	Db.Log.Warn("Disconnecting from database")
	defer Db.cancel()
	if err := Db.client.Disconnect(Db.ctx); err != nil {
		Db.Log.Panicf("Error trying close connection db: %s", err.Error())
	}
}
