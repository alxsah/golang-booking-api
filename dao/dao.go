package dao

import (
	"log"
	mgo "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  . "go-api/booking"
)

type BookingsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "bookings"
)

func (m *BookingsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *BookingsDAO) FindAll() ([]Booking, error) {
	var bookings []Booking
	err := db.C(COLLECTION).Find(bson.M{}).All(&bookings)
	return bookings, err
}

func (m *BookingsDAO) FindById(id string) (Booking, error) {
	var booking Booking
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&booking)
	return booking, err
}

func (m *BookingsDAO) Insert(booking Booking) error {
	err := db.C(COLLECTION).Insert(&booking)
	return err
}

func (m *BookingsDAO) Delete(booking Booking) error {
	err := db.C(COLLECTION).Remove(&booking)
	return err
}

func (m *BookingsDAO) Update(id string, booking Booking) error {
  err := db.C(COLLECTION).UpdateId(bson.ObjectIdHex(id), 
    bson.M{"$set": bson.M{"name": booking.Name, "date": booking.Date}})
	return err
}