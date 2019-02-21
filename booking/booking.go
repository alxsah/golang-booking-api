package booking

import (
  "gopkg.in/mgo.v2/bson"
)

type Booking struct {
  ID    bson.ObjectId   `bson:"_id" json:"id"`
  Name 	string 			    `bson:"name" json:"name"`
  Date	string			    `bson:"date" json:"date"`
}