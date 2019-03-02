package auth

import (
  "time"
  "fmt"
  "net/http"
  jwt "github.com/dgrijalva/jwt-go"
  . "github.com/alxsah/golang-booking-api/utils"
)

type Token struct {
  Token   string  `json:"token"`
}

var mySigningKey = []byte("mysecret")

func parseJWT(header string) (*jwt.Token, error) {
  var tokenStr string
  _, err := fmt.Sscanf(header, "Bearer %s", &tokenStr);
  if err != nil {
    return nil, err
  }
  token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, err
    }
    return mySigningKey, nil
  })
  return token, err
}

func GenerateJWT() (string, error) {
  token := jwt.New(jwt.SigningMethodHS256)
  claims := token.Claims.(jwt.MapClaims)

  claims["authorized"] = true
  claims["username"] = "alexs"
  claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

  tokenString, err := token.SignedString(mySigningKey)

  if err != nil {
    fmt.Errorf("Something went wrong: %s", err.Error())
    return "", err
  }
  return tokenString, nil
}

func IsAuthorized(
  endpoint func(http.ResponseWriter, *http.Request)) func(
    http.ResponseWriter,
    *http.Request,
  ) {
  return func (w http.ResponseWriter, r *http.Request) {
    if r.Header["Authorization"] != nil {
      token, err := parseJWT(r.Header["Authorization"][0])
      if err != nil {
        RespondWithError(w, http.StatusUnauthorized, "Invalid JWT")
      } else if token.Valid {
        endpoint(w, r)
      }
    } else {
      RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
    }
  }
}