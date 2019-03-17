package main

import (
	"encoding/json"
	. "github.com/alxsah/golang-booking-api/auth"
	. "github.com/alxsah/golang-booking-api/booking"
	. "github.com/alxsah/golang-booking-api/config"
	. "github.com/alxsah/golang-booking-api/dao"
	. "github.com/alxsah/golang-booking-api/user"
	. "github.com/alxsah/golang-booking-api/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

var config = Config{}
var dao = BookingsDAO{}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := dao.CreateUser(user); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}
	RespondWithJson(w, http.StatusCreated, map[string]string{"message": "Success"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	userFound, err := dao.ValidateUser(user)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	token, err := getToken(userFound.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate JWT")
		return
	}
	RespondWithJson(w, http.StatusAccepted, token)
}

func getToken(uid bson.ObjectId) (Token, error) {
	var token Token
	tokenString, err := GenerateJWT(uid)
	if err != nil {
		return token, err
	}
	token = Token{tokenString}
	return token, nil
}

func findBookingById(r *http.Request, uid string) (Booking, error) {
	params := mux.Vars(r)
	booking, err := dao.FindById(params["id"], uid)
	return booking, err
}

func getAllBookings(w http.ResponseWriter, r *http.Request, uid string) {
	bookings, err := dao.FindAll(uid)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, bookings)
}

func getBooking(w http.ResponseWriter, r *http.Request, uid string) {
	booking, err := findBookingById(r, uid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Booking ID")
		return
	}
	RespondWithJson(w, http.StatusOK, booking)
}

func createBooking(w http.ResponseWriter, r *http.Request, uid string) {
	defer r.Body.Close()
	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	booking.ID = bson.NewObjectId()
	booking.UID = bson.ObjectIdHex(uid)
	if err := dao.Insert(booking); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusCreated, booking)
}

func deleteBooking(w http.ResponseWriter, r *http.Request, uid string) {
	booking, err := findBookingById(r, uid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Booking ID")
		return
	}
	if err := dao.Delete(booking, uid); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func updateBooking(w http.ResponseWriter, r *http.Request, uid string) {
	defer r.Body.Close()
	var newBooking Booking
	id := mux.Vars(r)["id"]
	_, err := findBookingById(r, uid)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Booking ID")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&newBooking); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(id, newBooking); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func init() {
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/bookings", IsAuthorized(getAllBookings)).Methods("GET")
	r.HandleFunc("/bookings", IsAuthorized(createBooking)).Methods("POST")
	r.HandleFunc("/bookings/{id}", IsAuthorized(updateBooking)).Methods("PUT")
	r.HandleFunc("/bookings/{id}", IsAuthorized(getBooking)).Methods("GET")
	r.HandleFunc("/bookings/{id}", IsAuthorized(deleteBooking)).Methods("DELETE")

	if err := http.ListenAndServe(":3001",
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(r)); err != nil {
		log.Fatal(err)
	}
}
