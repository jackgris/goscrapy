package database

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wholesalers struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Login      string             `bson:"login"`
	User       string             `bson:"user"`
	Pass       string             `bson:"pass"`
	Searchpage string             `bson:"serchpage"`
	Name       string             `bson:"name"`
}

func (db *MongoDb) InsertWholesaer(w Wholesalers) error {

	collection := db.client.Database(db.name).Collection("provider")
	whosaler := db.GetWhosalerById(w)
	var err error
	// If not exist inserted in other case update data
	if whosaler.Name != w.Name {
		_, err = collection.InsertOne(db.ctx, w)
	} else {
		filter := bson.M{"_id": w.Id}
		_, err = collection.UpdateOne(db.ctx, filter, w)
	}
	return err
}

// Reading from database, a product identified with his ID
func (m *MongoDb) GetWhosalerById(w Wholesalers) Wholesalers {

	collection := db.client.Database(m.name).Collection("provider")
	r := collection.FindOne(db.ctx, bson.M{"_id": w.Id})
	wholesaler := Wholesalers{}
	err := r.Decode(&wholesaler)

	if err != nil {
		fmt.Println("Error GetWhosalerById cursor: ", err)
	}

	return wholesaler
}

// Reading from database and returning all wholesalers
func (m *MongoDb) FindWholesalers() []Wholesalers {

	whosalers := []Wholesalers{}
	collection := db.client.Database(db.name).Collection("provider")
	filter := bson.D{}
	cur, err := collection.Find(m.ctx, filter)
	if err != nil {
		fmt.Println("Error FindWhosalers getting cursor: ", err)
		return whosalers
	}
	defer cur.Close(db.ctx)

	for cur.Next(db.ctx) {
		var result Wholesalers
		err := cur.Decode(&result)
		if err != nil {
			fmt.Println("Error FindWhosalers decode bson: ", err)
		}

		whosalers = append(whosalers, result)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error FindWhosalers cursor: ", err)
	}
	return whosalers
}

// Delete a product from DB, if can't do this, return an error. And will print when not found matchs
func (m *MongoDb) DeleteWhosaler(w Wholesalers) error {

	collection := db.client.Database(db.name).Collection("provider")
	result, err := collection.DeleteOne(db.ctx, bson.M{"_id": w.Id})
	if result.DeletedCount == 0 {
		return errors.New("Delete wholesaler not found match")
	}
	if result.DeletedCount > 1 {
		return errors.New("Delete wholesaler found many matchs")
	}
	return err
}
