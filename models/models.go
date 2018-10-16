package models

import "github.com/globalsign/mgo/bson"

type Planet struct {
	ID          bson.ObjectId `json:"id" bson:"_id" valid:"-"`
	Climate     string        `json:"climate" bson:"climate" valid:"alphanum"`
	Name        string        `json:"name" bson:"name" valid:"alphanum"`
	Terrain     string        `json:"terrain" bson:"terrain" valid:"alphanum"`
	Appearances int           `json:"appearances" bson:"appearances" valid:"-"`
}