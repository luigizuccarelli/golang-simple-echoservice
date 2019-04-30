package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAllMiddleware(t *testing.T) {
	var req *http.Request
	var token string = ""
	var response Response

	// create anonymous struct
	tests := []struct {
		Name     string
		Method   string
		Url      string
		Payload  string
		Handler  string
		FileName string
		want     int
		errorMsg string
	}{
		{
			"Test [Isalive] should pass",
			"GET", "/v1/sys/info/isalive",
			"",
			"IsAlive",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"Test [MiddlewareLogin] should pass",
			"POST",
			"/v1/login",
			"{\"username\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMxNyA5QVQ=\"}",
			"MiddlewareLogin",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"Test [MiddlewareData] should pass",
			"POST",
			"/v1/alldata",
			"{\"apitoken\": \"\"}",
			"MiddlewareData",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"Test [MiddlewareCustomerNumberData] should pass",
			"GET",
			"/v2/postaladdress/customernumber/000002096234",
			"",
			"MiddlewareCustomerNumberData",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"Test [MiddlewareLogin] should fail",
			"POST",
			"/v1/login",
			"{\"user\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"pass\":\"TFMxNyA5QVQ=\"}",
			"MiddlewareLogin",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"Test [MiddlewareData] should fail",
			"POST",
			"/v1/alldata",
			"{\"apitoken\": \"\"}",
			"MiddlewareLogin",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"Test [MiddlewareCustomerNumberData] should fail",
			"GET",
			"/v2/postaladdress/customernumber/00000",
			"",
			"MiddlewareCustomerNumberData",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
	}

	for _, tt := range tests {
		fmt.Printf("Executing test : %s \n", tt.Name)
		if tt.Payload == "" {
			req, _ = http.NewRequest(tt.Method, tt.Url, nil)
		} else {
			if strings.Contains(tt.Payload, "apitoken") {
				tt.Payload = "{\"apitoken\": \"" + token + "\"}"
			}
			req, _ = http.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Payload)))
		}

		connectors = NewTestClients(tt.FileName, tt.want)

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		switch tt.Handler {
		case "IsAlive":
			handler := http.HandlerFunc(IsAlive)
			handler.ServeHTTP(rr, req)
		case "MiddlewareLogin":
			handler := http.HandlerFunc(MiddlewareLogin)
			handler.ServeHTTP(rr, req)
		case "MiddlewareData":
			handler := http.HandlerFunc(MiddlewareData)
			handler.ServeHTTP(rr, req)
		case "MiddlewareCustomerNumberData":
			router := mux.NewRouter()
			router.HandleFunc("/v2/postaladdress/customernumber/{customernumber}", MiddlewareCustomerNumberData)
			router.ServeHTTP(rr, req)
		}
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		fmt.Println(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		json.Unmarshal(body, &response)
		token = response.Payload.MetaInfo
		if rr.Code != tt.want {
			t.Errorf(tt.errorMsg, rr.Code, http.StatusOK)
		}
	}
}
