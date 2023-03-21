package database

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wholesalers struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Login        string             `json:"login" bson:"login"`
	User         string             `json:"user" bson:"user"`
	Pass         string             `json:"pass" bson:"pass"`
	Searchpage   string             `json:"searchpage" bson:"searchpage"`
	Name         string             `json:"name" bson:"name"`
	EndPhrase    string             `json:"endphrase"`
	EndPhraseDiv string             `json:"endphrasediv"`
}

func (db *MongoDb) InsertWholesaer(w Wholesalers) error {

	collection := db.client.Database(db.name).Collection("provider")
	whosaler := db.GetWhosalerName(w)
	var err error
	// If not exist inserted in other case update data
	if whosaler.Name != w.Name {
		_, err = collection.InsertOne(db.ctx, w)
	} else {
		w.Id = whosaler.Id
		filter := bson.M{"_id": w.Id}
		update := bson.M{"$set": bson.M{
			"user":         w.User,
			"pass":         w.Pass,
			"login":        w.Login,
			"searchpage":   w.Searchpage,
			"endphrase":    w.EndPhrase,
			"endphrasediv": w.EndPhraseDiv,
		}}
		_, err = collection.UpdateOne(db.ctx, filter, update)
	}
	return err
}

// Reading from database, a product identified with his ID
func (m *MongoDb) GetWhosalerName(w Wholesalers) Wholesalers {

	collection := Db.client.Database(m.name).Collection("provider")
	r := collection.FindOne(Db.ctx, bson.M{"name": w.Name})
	wholesaler := Wholesalers{}
	err := r.Decode(&wholesaler)

	if err != nil {
		m.Log.Info("GetWhosalerName cursor: ", err)
	}

	return wholesaler
}

// Reading from database, a product identified with his ID
func (m *MongoDb) GetWhosalerById(w Wholesalers) Wholesalers {

	collection := Db.client.Database(m.name).Collection("provider")
	r := collection.FindOne(Db.ctx, bson.M{"_id": w.Id})
	wholesaler := Wholesalers{}
	err := r.Decode(&wholesaler)

	if err != nil {
		m.Log.Info("GetWhosalerById cursor: ", err)
	}

	return wholesaler
}

// Reading from database and returning all wholesalers
func (m *MongoDb) FindWholesalers() []Wholesalers {

	whosalers := []Wholesalers{}
	collection := Db.client.Database(Db.name).Collection("provider")
	filter := bson.D{}
	cur, err := collection.Find(m.ctx, filter)
	if err != nil {
		m.Log.Info("FindWhosalers getting cursor: ", err)
		return whosalers
	}
	defer cur.Close(Db.ctx)

	for cur.Next(Db.ctx) {
		var result Wholesalers
		err := cur.Decode(&result)
		if err != nil {
			m.Log.Info("FindWhosalers decode bson: ", err)
		}

		whosalers = append(whosalers, result)
	}

	if err := cur.Err(); err != nil {
		m.Log.Info("FindWhosalers cursor: ", err)
	}
	return whosalers
}

// Delete a product from DB, if can't do this, return an error. And will print when not found matchs
func (m *MongoDb) DeleteWhosaler(w Wholesalers) error {

	collection := Db.client.Database(Db.name).Collection("provider")
	result, err := collection.DeleteOne(Db.ctx, bson.M{"_id": w.Id})
	if result.DeletedCount == 0 {
		return errors.New("Delete wholesaler not found match")
	}
	if result.DeletedCount > 1 {
		return errors.New("Delete wholesaler found many matchs")
	}
	return err
}
