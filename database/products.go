package database

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Data product needed, not mather what wholesaler get data
type Product struct {
	Id_         primitive.ObjectID `bson:"_id,omitempty"`
	Id          string             `bson:"id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Image       string
	Description string
	Price       []Value `bson:"prices,omitempty"`
	Stock       string
	Wholesaler  string
}

type Value struct {
	Price float64   `bson:"price,omitempty"`
	Date  time.Time `bson:"date,omitempty"`
}

// Inserting product in the database, only when is a new product or when price change,
// is can't do this, return error
func (db *MongoDb) Create(p Product) error {

	collection := db.client.Database(db.name).Collection("productos")
	product := db.ReadById(p)
	var err error
	// If product doesn't exist, insert a new one, but if exist only update prices
	if product.Id == "" {
		_, err = collection.InsertOne(db.ctx, p)
	} else {
		// update only if last price are different
		if product.Price[len(product.Price)-1].Price != p.Price[len(p.Price)-1].Price {
			products := product.Price
			products = append(products, p.Price[len(p.Price)-1])
			product.Price = products
			err = db.updateProduct(product)
		}
	}
	return err
}

func (db *MongoDb) updateProduct(p Product) error {
	collection := db.client.Database(db.name).Collection("productos")
	filter := bson.M{"_id": p.Id_}
	update := bson.M{"$set": bson.M{
		"id":          p.Id,
		"name":        p.Name,
		"image":       p.Image,
		"description": p.Description,
		"prices":      p.Price,
		"stock":       p.Stock,
		"wholesaler":  p.Wholesaler,
	}}
	_, err := collection.UpdateOne(db.ctx, filter, update)
	return err
}

// Reading from database, a product identified with his ID
func (db *MongoDb) ReadById(p Product) Product {

	collection := db.client.Database(db.name).Collection("productos")
	r := collection.FindOne(db.ctx, bson.M{"id": p.Id})
	var product Product
	// If can't decode data, means that doesn't exist that product.
	// So can ignored the error because we'll return an empty product.
	_ = r.Decode(&product)

	return product
}

// Reading from database, a product identified with his ID
func (db *MongoDb) ReadByMongoId(p Product) Product {

	collection := db.client.Database(db.name).Collection("productos")
	r := collection.FindOne(db.ctx, bson.M{"_id": p.Id_})
	product := Product{}
	err := r.Decode(&product)

	if err != nil {
		db.Log.Info("ReadByMongoId cursor: ", err)
	}

	return product
}

// Reading from database and returning all products from one wholesaler
func (db *MongoDb) ReadByWholesalers(name string) []Product {

	products := []Product{}
	collection := db.client.Database(db.name).Collection("productos")
	filter := bson.M{"wholesaler": name}
	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		db.Log.Info("ReadByWholesaler getting cursor: ", err)
		return products
	}
	defer cur.Close(db.ctx)

	for cur.Next(db.ctx) {
		var result Product
		err := cur.Decode(&result)
		if err != nil {
			db.Log.Info("ReadByWholesalers decode bson: ", err)
		}

		products = append(products, result)
	}

	if err := cur.Err(); err != nil {
		db.Log.Info("ReadByWholersalers cursor: ", err)
	}
	return products
}

// Delete a product from DB, if can't do this, return an error. And will print when not found matchs
func (db *MongoDb) Delete(p Product) error {

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
