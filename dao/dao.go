package dao

import (
  "log"
  mgo "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  . "github.com/alxsah/golang-booking-api/booking"
  . "github.com/alxsah/golang-booking-api/user"
)

type BookingsDAO struct {
  Server   string
  Database string
}

var db *mgo.Database

const (
  COLLECTION_BOOKINGS = "bookings"
  COLLECTION_USERS = "users"
)

func (m *BookingsDAO) Connect() {
  session, err := mgo.Dial(m.Server)
  if err != nil {
    log.Fatal(err)
  }
  db = session.DB(m.Database)
}

func (m *BookingsDAO) FindAll(uid string) ([]Booking, error) {
  bookings := make([]Booking, 0)
  err := db.C(COLLECTION_BOOKINGS).Find(bson.M{"uid": bson.ObjectIdHex(uid)}).All(&bookings)
  return bookings, err
}

func (m *BookingsDAO) FindById(id string, uid string) (Booking, error) {
  var booking Booking
  err := db.C(COLLECTION_BOOKINGS).Find(bson.M{
    "_id": bson.ObjectIdHex(id),
    "uid": bson.ObjectIdHex(uid),
  }).One(&booking)
  return booking, err
}

func (m *BookingsDAO) Insert(booking Booking) error {
  booking.ID = bson.NewObjectId()
  err := db.C(COLLECTION_BOOKINGS).Insert(&booking)
  return err
}

func (m *BookingsDAO) Delete(booking Booking, uid string) error {
  err := db.C(COLLECTION_BOOKINGS).Remove(bson.M{
    "_id": booking.ID,
    "uid": bson.ObjectIdHex(uid),
  })
  return err
}

func (m *BookingsDAO) Update(id string, booking Booking) error {
  err := db.C(COLLECTION_BOOKINGS).UpdateId(bson.ObjectIdHex(id), 
    bson.M{"$set": bson.M{
      "name": booking.Name, 
      "date": booking.Date,
      "location": booking.Location,
    }})
  return err
}

func (m *BookingsDAO) CreateUser(user User) error {
  user.ID = bson.NewObjectId()
  var userFound User
  err := db.C(COLLECTION_USERS).Find(bson.M{"username": user.Username}).One(&userFound)
  if err != nil && err != mgo.ErrNotFound {
    return err
  }
  err = db.C(COLLECTION_USERS).Insert(&user)
  return err
}

func (m *BookingsDAO) ValidateUser(user User) (User, error) {
  var userFound User
  err := db.C(COLLECTION_USERS).Find(bson.M{
    "username": user.Username,
    "password": user.Password,
  }).One(&userFound)
  return userFound, err
}