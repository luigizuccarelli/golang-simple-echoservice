package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-authinterface/pkg/connectors"
	"github.com/microlib/simple"
)

func TestAllMiddleware(t *testing.T) {

	logger := &simple.Logger{Level: "info"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("AuthHandler : should fail (no token)", func(t *testing.T) {
		var STATUS int = 403
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/verify", nil)
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			AuthHandler(w, req, con)
		})
		//req.Header.Set("Authorization", ": Bearer dsfdsfdfdsfdsfd")

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "AuthHandler", rr.Code, STATUS))
		}
	})

	t.Run("AuthHandler : should fail (bad request)", func(t *testing.T) {
		var STATUS int = 400
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/verify", nil)
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			AuthHandler(w, req, con)
		})

		req.Header.Set("Authorization", ": Bearer dsfdsfdfdsfdsfd")

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "AuthHandler", rr.Code, STATUS))
		}
	})

	t.Run("AuthHandler : should fail (token invalid)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("JWT_SECRETKEY", "uraidiot")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/verify", nil)
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			AuthHandler(w, req, con)
		})

		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.oYrIXuRzmxX0DpKDbPhzDo0UtasgmWWVCvjGHYdXS74")

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "AuthHandler", rr.Code, STATUS))
		}
	})

	t.Run("AuthHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("JWT_SECRETKEY", "uratool")

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/verify", nil)
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			AuthHandler(w, req, con)
		})

		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.oYrIXuRzmxX0DpKDbPhzDo0UtasgmWWVCvjGHYdXS74")

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "AuthHandler", rr.Code, STATUS))
		}
	})

}

/*
	// create anonymous struct
	tests := []struct {
		Name     string
		Method   string
		Url      string
		Payload  string
		Handler  string
		FileName string
		Want     int
		ErrorMsg string
	}{
		{
			"[TEST] IsAlive should pass",
			"GET", "/v1/sys/info/isalive",
			"",
			"IsAlive",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] Method [OPTIONS] should pass",
			"OPTIONS", "/v1/sys/info/isalive",
			"",
			"IsAlive",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] MiddlewareAuth should fail",
			"GET",
			"/api/v1/verify",
			"{\"user\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"pass\":\"TFMxNyA5QVQ=\"}",
			"MiddlewareAuth",
			"tests/payload-example.json",
			http.StatusBadRequest,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] MiddlewareAuth no token should fail",
			"GET",
			"/api/v1/verify",
			"{\"user\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"pass\":\"TFMxNyA5QVQ=\"}",
			"MiddlewareAuthNoToken",
			"tests/payload-example.json",
			http.StatusForbidden,
			"Handler returned wrong status code got %d want %d",
		},
	}

	for _, tt := range tests {
		logger.Info(fmt.Sprintf("%s : \n", tt.Name))
		if tt.Payload == "" {
			req, _ = http.NewRequest(tt.Method, tt.Url, nil)
		} else {
			if strings.Contains(tt.Payload, "apitoken") {
				tt.Payload = "{\"apitoken\": \"" + token + "\"}"
			}
			req, _ = http.NewRequest(tt.Method, tt.Url, bytes.NewBuffer([]byte(tt.Payload)))
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		switch tt.Handler {
		case "IsAlive":
			handler := http.HandlerFunc(IsAlive)
			handler.ServeHTTP(rr, req)
		case "MiddlewareAuth":
			handler := http.HandlerFunc(MiddlewareAuth)
			req.Header.Set("Authorization", " Bearer rewrtewrewrwerwerwerwerewr")
			handler.ServeHTTP(rr, req)
		case "MiddlewareAuthNoToken":
			handler := http.HandlerFunc(MiddlewareAuth)
			handler.ServeHTTP(rr, req)
		}
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		_, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		if rr.Code != tt.Want {
			t.Errorf(tt.ErrorMsg, rr.Code, tt.Want)
		}
	}
}*/
