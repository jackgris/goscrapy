package database

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Data product needed, not mather what wholesaler get data
type Product struct {
	Id_         primitive.ObjectID `bson:"_id,omitempty"`
	Id          string
	Name        string
	Image       string
	Description string
	Price       string
	Stock       string
	Wholesaler  string
}

// My database struct
type MongoDb struct {
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
	name   string
}

// vars needed for only created one instance for my access to the database
var (
	once sync.Once
	db   *MongoDb
)

// Connecting to the database, only one instance will be create, to connect with mongo database, we need
// the URI where are the DB, the user name and password, and will return the instance with an active connection
// or an error
func Connect(dburi, dbuser, dbpass, name string) (*MongoDb, error) {

	var err error
	once.Do(func() {
		db = new(MongoDb)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		db.cancel = cancel
		db.ctx = ctx
		option := options.Client().ApplyURI(dburi)
		credentials := options.Credential{Username: dbuser, Password: dbpass}
		option.SetAuth(credentials)

		db.client, err = mongo.Connect(db.ctx, option)
		if err != nil {
			panic("Can't connect with db" + err.Error())
		}
		err = db.client.Ping(db.ctx, readpref.Primary())

	})

	db.name = name

	return db, err
}

// Closing the database connection, correctly
func Disconnect() {
	fmt.Println("Disconnecting from database")
	defer db.cancel()
	if err := db.client.Disconnect(db.ctx); err != nil {
		panic("Error trying close connection db: " + err.Error())
	}
}

// Inserting product in the database, only when is a new product or when price change,
// is can't do this, return error
func (db *MongoDb) Create(p Product) error {

	collection := db.client.Database(db.name).Collection("productos")
	product := db.ReadById(p)
	var err error
	if product.Id != "" {
		if product.Price != p.Price {
			_, err = collection.InsertOne(db.ctx, p)
		}
	} else {
		_, err = collection.InsertOne(db.ctx, p)
	}
	return err
}

// Reading from database, a product identified with his ID
func (m *MongoDb) ReadById(p Product) Product {

	collection := db.client.Database(m.name).Collection("productos")
	r := collection.FindOne(db.ctx, bson.M{"id": p.Id})
	product := Product{}
	err := r.Decode(&product)

	if err != nil {
		fmt.Println("Error ReadById cursor: ", err)
	}

	return product
}

// Reading from database, a product identified with his ID
func (m *MongoDb) ReadByMongoId(p Product) Product {

	collection := db.client.Database(m.name).Collection("productos")
	r := collection.FindOne(db.ctx, bson.M{"_id": p.Id_})
	product := Product{}
	err := r.Decode(&product)

	if err != nil {
		fmt.Println("Error ReadByMongoId cursor: ", err)
	}

	return product
}

// Reading from database and returning all products from one wholesaler
func (m *MongoDb) ReadByWholesalers(name string) []Product {

	products := []Product{}
	collection := db.client.Database(db.name).Collection("productos")
	filter := bson.M{"wholesaler": name}
	cur, err := collection.Find(m.ctx, filter)
	if err != nil {
		fmt.Println("Error ReadByWholesaler getting cursor: ", err)
		return products
	}
	defer cur.Close(db.ctx)

	for cur.Next(db.ctx) {
		var result Product
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println("Error ReadByWholesalers decode bson: ", err)
		}

		products = append(products, result)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error ReadByWholersalers cursor: ", err)
	}
	return products
}

// Delete a product from DB, if can't do this, return an error. And will print when not found matchs
func (m *MongoDb) Delete(p Product) error {

	collection := db.client.Database(db.name).Collection("productos")
	result, err := collection.DeleteOne(db.ctx, bson.M{"_id": p.Id_})
	if result.DeletedCount == 0 {
		return errors.New("Delete not found match")
	}
	if result.DeletedCount > 1 {
		return errors.New("Delete found many matchs")
	}
	return err
}
