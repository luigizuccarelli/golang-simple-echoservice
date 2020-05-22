package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-authinterface/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-authinterface/pkg/schema"
	"github.com/dgrijalva/jwt-go"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

type Claims struct {
	jwt.StandardClaims
}

func AuthHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var response *schema.Response

	token := r.Header.Get(strings.ToLower("Authorization"))
	conn.Trace("Auth Header : %s", token)

	if token == "" {
		w.WriteHeader(http.StatusForbidden)
		response = &schema.Response{Code: 403, StatusCode: "403", Status: "ERROR", Message: "Forbidden", Payload: schema.SchemaInterface{}}
	} else {

		// Remove Bearer
		tknStr := strings.Trim(token[7:], " ")
		conn.Info(fmt.Sprintf("Token (trimmed) : %s", tknStr))
		addHeaders(w, r)

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		secret := os.Getenv("JWT_SECRETKEY")
		conn.Trace(fmt.Sprintf("JWT SECRET : %s", secret))

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
				w.WriteHeader(http.StatusForbidden)
				response = &schema.Response{Code: 403, StatusCode: "403", Status: "ERROR", Message: "Forbidden", Payload: schema.SchemaInterface{}}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				response = &schema.Response{Code: 400, StatusCode: "400", Status: "ERROR", Message: "Bad Request", Payload: schema.SchemaInterface{}}
			}
		} else {
			if !tkn.Valid {
				w.WriteHeader(http.StatusForbidden)
				response = &schema.Response{Code: 403, StatusCode: "403", Status: "ERROR", Message: "Forbidden", Payload: schema.SchemaInterface{}}
			} else {
				response = &schema.Response{Code: 200, StatusCode: "200", Status: "OK", Message: "Data uploaded succesfully", Payload: schema.SchemaInterface{}}
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	conn.Debug(fmt.Sprintf("AuthHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \"AuthInterface\" }")
	return
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
