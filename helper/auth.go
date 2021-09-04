package helper

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/joho/godotenv"
	//    "github.com/gorilla/mux"
	//"github.com/mitchellh/mapstructure"
)

var mySigningKey []byte

func getSecret() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mySigningKey = []byte(os.Getenv("SECRET_KEY"))
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	getSecret()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return mySigningKey, nil
				})
				if error != nil {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				if token.Valid {
					context.Set(r, "decoded", token.Claims)
					next(w, r)
				} else {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	})
}
