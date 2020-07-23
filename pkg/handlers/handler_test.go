package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-simple-oc4service/pkg/connectors"
	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Inject (force) readAll test error")
}

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

	t.Run("EchoHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		requestPayload := `{  "id": "12345566", "message": "Hello World - from Luigi" }`
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/echo", bytes.NewBuffer([]byte(requestPayload)))
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			EchoHandler(w, req, con)
		})
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

	t.Run("EchoHandler : should fail (force readall error)", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/echo", errReader(0))
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			EchoHandler(w, req, con)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EchoHandler", rr.Code, STATUS))
		}
	})

	t.Run("EchoHandler : should fail (bad json)", func(t *testing.T) {
		var STATUS int = 500
		requestPayload := `{  "id": "12345566", " }`
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/echo", bytes.NewBuffer([]byte(requestPayload)))
		con := connectors.NewTestConnectors(logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			EchoHandler(w, req, con)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EchoHandler", rr.Code, STATUS))
		}
	})

}
