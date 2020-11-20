package database

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database interface {
	Create(p Product) error
	ReadById(p Product) []Product
	ReadByWholesalers(name string) []Product
	Delete(p Product) error
}

type Product struct {
	Id          string
	Name        string
	Image       string
	Description string
	Price       string
	Stock       string
	Wholesaler  string
}

type MongoDb struct {
	client *mongo.Client
	ctx    context.Context
}

var (
	once sync.Once
	db   *MongoDb
)

func Connect(dburi, dbuser, dbpass string, c context.Context) (*MongoDb, error) {

	var err error
	once.Do(func() {
		db = new(MongoDb)
		// FIXME should use context with timeout
		//ctx, cancel := context.WithTimeout(c, 10*time.Second)
		//defer cancel()
		db.ctx = c
		option := options.Client().ApplyURI(dburi)
		credentials := options.Credential{Username: dbuser, Password: dbpass}
		option.SetAuth(credentials)

		db.client, err = mongo.Connect(db.ctx, option)
		if err != nil {
			panic("Can't connect with db" + err.Error())
		}
		err = db.client.Ping(db.ctx, readpref.Primary())

	})
	return db, err
}

func Disconnect() {
	fmt.Println("Disconnecting from database")
	if err := db.client.Disconnect(db.ctx); err != nil {
		panic("Error trying close connection db: " + err.Error())
	}
}

func (db *MongoDb) Create(p Product) error {

	collection := db.client.Database("mayorista").Collection("productos")
	products := db.ReadById(p)
	// FIXME Should delete this print when test the function
	var err error
	if len(products) > 0 {
		if products[len(products)-1].Price != p.Price {
			_, err = collection.InsertOne(db.ctx, p)
		}
	} else {
		_, err = collection.InsertOne(db.ctx, p)
	}
	return err
}

func (m *MongoDb) ReadById(p Product) []Product {

	products := []Product{}
	collection := db.client.Database("mayorista").Collection("productos")
	cur, err := collection.Find(db.ctx, bson.M{"id": p.Id})
	if err != nil {
		fmt.Println("Error ReadById getting cursor: ", err)
		return products
	}
	defer cur.Close(db.ctx)
	for cur.Next(db.ctx) {
		var result Product
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println("Error ReadById decode bson: ", err)
		}

		products = append(products, result)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error ReadById cursor: ", err)
	}

	return products
}

func (m *MongoDb) ReadByWholesalers(name string) []Product {

	products := []Product{}
	collection := db.client.Database("mayorista").Collection("productos")
	cur, err := collection.Find(m.ctx, Product{Wholesaler: name})
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

func (m *MongoDb) Delete(p Product) error {

	collection := db.client.Database("mayorista").Collection("productos")
	result, err := collection.DeleteOne(db.ctx, p)
	if result.DeletedCount == 0 {
		fmt.Println("Delete not found match")
	}
	if result.DeletedCount == 1 {
		fmt.Println("Delete found many matchs")
	}
	return err
}
