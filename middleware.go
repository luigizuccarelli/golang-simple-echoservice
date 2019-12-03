package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

type Claims struct {
	jwt.StandardClaims
}

func MiddlewareAuth(w http.ResponseWriter, r *http.Request) {
	var response Response

	token := r.Header.Get(strings.ToLower("Authorization"))

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		response = Response{Code: 403, StatusCode: "403", Status: "ERROR", Message: "Forbidden", Payload: SchemaInterface{}}
	} else {

		// Remove Bearer
		tknStr := strings.Trim(token[7:], " ")
		logger.Debug(fmt.Sprintf("Token : %s", tknStr))
		addHeaders(w, r)
		handleOptions(w, r)

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		secret := os.Getenv("JWT_SECRETKEY")
		logger.Trace(fmt.Sprintf("JWT SECRET : %s", secret))

		var jwtKey = []byte(secret)

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err.Error() == jwt.ErrSignatureInvalid.Error() {
				w.WriteHeader(http.StatusUnauthorized)
				response = Response{Code: 403, StatusCode: "403", Status: "ERROR", Message: "Forbidden", Payload: SchemaInterface{}}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				response = Response{Code: 400, StatusCode: "400", Status: "ERROR", Message: "Bad Request", Payload: SchemaInterface{}}
			}
		} else {
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				response = Response{Code: 403, StatusCode: "403", Status: "ERROR", Message: "Forbidden", Payload: SchemaInterface{}}
			} else {
				response = Response{Code: 200, StatusCode: "200", Status: "OK", Message: "Data uploaded succesfully", Payload: SchemaInterface{}}
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("AuthHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \"Auth\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	var request []string
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	logger.Trace(fmt.Sprintf("Headers : %s", request))

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

// simple options handler
func handleOptions(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "OPTIONS" {
		return true
	}
	return false
}

// simple error handler
func handleError(w http.ResponseWriter, msg string) Response {
	w.WriteHeader(http.StatusInternalServerError)
	r := Response{Code: 500, StatusCode: "500", Status: "ERROR", Message: msg}
	return r
}
