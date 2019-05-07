package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// MiddlewareLogin a http response and request wrapper for database insert
// It takes a both response and request objects and returns void
func MiddlewareLogin(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	addHeaders(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareLogin) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	apitoken, err := connectors.LoginData(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Login error (MiddlewareLogin) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		payload = SchemaInterface{MetaInfo: apitoken}
		response = Response{StatusCode: "200", Status: "OK", Message: "MW data successfully retrieved", Payload: payload}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// MiddlewareData a http response and request wrapper for database update
// It takes a both response and request objects and returns void
func MiddlewareData(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface

	addHeaders(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Could not read body data (MiddlewareData) " + err.Error(), Payload: payload}
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := connectors.AllData(body)
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Subscriptions data read (MiddlewareSubscriptions) " + err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		e := json.Unmarshal(data, &payload)
		if e != nil {
			response = Response{StatusCode: "500", Status: "ERROR", Message: "Subscriptions unmarshal error (MiddlewareSubscriptions) " + e.Error()}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response = Response{StatusCode: "200", Status: "OK", Message: "MW data read successfully", Payload: payload}
		}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

// It takes a both response and request objects and returns void
func MiddlewareCustomerNumberData(w http.ResponseWriter, r *http.Request) {

	var response Response
	var payload SchemaInterface
	var vars = mux.Vars(r)

	addHeaders(w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
		return
	}

	logger.Info(fmt.Sprintf("In function MiddlewareCustomerNumberData %v", vars))

	data, err := connectors.AllDataByCustomerNumber(vars["customernumber"])
	if err != nil {
		response = Response{StatusCode: "500", Status: "ERROR", Message: "Accounts data read (MiddlewareCustomerNumberData) " + err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		e := json.Unmarshal(data, &payload)
		if e != nil {
			response = Response{StatusCode: "500", Status: "ERROR", Message: "Account unmarshal error (MiddlewareCustomerNumberData) " + e.Error()}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response = Response{StatusCode: "200", Status: "OK", Message: "MW Account data read successfully", Payload: payload}
		}
	}

	b, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	logger.Debug(fmt.Sprintf("used to mask cc %v", r))
	logger.Trace(fmt.Sprintf("config data  %v", config))
	fmt.Fprintf(w, "{\"isalive\": true , \"version\": \"1.0.0\"}")
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
