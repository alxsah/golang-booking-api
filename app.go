package main
import (
  "encoding/json"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2/bson"
  . "github.com/alxsah/golang-booking-api/dao"
  . "github.com/alxsah/golang-booking-api/config"
  . "github.com/alxsah/golang-booking-api/booking"
  . "github.com/alxsah/golang-booking-api/auth"
  . "github.com/alxsah/golang-booking-api/utils"
)

var config = Config{}
var dao = BookingsDAO{}

func getToken(w http.ResponseWriter, r *http.Request) {
  tokenString, err := GenerateJWT()
  if err != nil {
    RespondWithError(w, http.StatusInternalServerError, "Failed to generate JWT")
    return
  }
  tokenObj := Token { tokenString }
  RespondWithJson(w, http.StatusOK, tokenObj)
}

func findBookingById(r *http.Request) (Booking, error) {
  params := mux.Vars(r)
  booking, err := dao.FindById(params["id"])
  return booking, err
}

func getAllBookings(w http.ResponseWriter, r *http.Request) {
  bookings, err := dao.FindAll()
  if err != nil {
    RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }
  RespondWithJson(w, http.StatusOK, bookings)
}

func getBooking(w http.ResponseWriter, r *http.Request) {
  booking, err := findBookingById(r)
  if err != nil {
    RespondWithError(w, http.StatusBadRequest, "Invalid Booking ID")
    return
  }
  RespondWithJson(w, http.StatusOK, booking)
}

func createBooking(w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()
  var booking Booking
  if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
    RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
  booking.ID = bson.NewObjectId()
  if err := dao.Insert(booking); err != nil {
    RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }
  RespondWithJson(w, http.StatusCreated, booking)
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
  booking, err := findBookingById(r);
  if err != nil {
    RespondWithError(w, http.StatusBadRequest, "Invalid Booking ID")
    return
  }
  if err := dao.Delete(booking); err != nil {
    RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }
  RespondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func updateBooking(w http.ResponseWriter, r *http.Request) {
  defer r.Body.Close()
  var newBooking Booking
  id := mux.Vars(r)["id"]
  _, err := findBookingById(r);
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

  r.HandleFunc("/get-token", getToken).Methods("GET")
  r.HandleFunc("/bookings", IsAuthorized(getAllBookings)).Methods("GET")
  r.HandleFunc("/bookings", IsAuthorized(createBooking)).Methods("POST")
  r.HandleFunc("/bookings/{id}", IsAuthorized(updateBooking)).Methods("PUT")
  r.HandleFunc("/bookings/{id}", IsAuthorized(getBooking)).Methods("GET")
  r.HandleFunc("/bookings/{id}", IsAuthorized(deleteBooking)).Methods("DELETE")

  if err := http.ListenAndServe(":3001", r); err != nil {
    log.Fatal(err)
  }
}

