package auth

import (
	"errors"
	"fmt"
	. "github.com/alxsah/golang-booking-api/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

var mySigningKey = []byte("mysecret")

func parseJWT(header string) (*jwt.Token, error) {
	var tokenStr string
	_, err := fmt.Sscanf(header, "Bearer %s", &tokenStr)
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

func GenerateJWT(uid bson.ObjectId) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["uid"] = uid
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(
	endpoint func(http.ResponseWriter, *http.Request, string)) func(
	http.ResponseWriter,
	*http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			token, err := parseJWT(r.Header["Authorization"][0])
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Invalid JWT")
			} else if token.Valid {
				uid, err := GetUIDFromToken(r.Header["Authorization"][0])
				if err != nil {
					RespondWithError(w, http.StatusUnauthorized, "Failed to retrieve UID")
				}
				endpoint(w, r, uid)
			}
		} else {
			RespondWithError(w, http.StatusUnauthorized, "Not Authorized")
		}
	}
}

func GetUIDFromToken(authHeader string) (string, error) {
	var tokenStr string
	claims := jwt.MapClaims{}
	if authHeader != "" {
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
		if err != nil {
			return "", err
		}
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
		if err != nil {
			return "", err
		}
		return claims["uid"].(string), nil
	}
	return "", errors.New("Auth header is empty")
}
