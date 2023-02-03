package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Wholesalers struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Login      string
	User       string
	Pass       string
	Searchpage string
	Name       string
}
