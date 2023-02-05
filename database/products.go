package database

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
