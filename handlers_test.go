package main

import (
	"bytes"
	"fmt"
	"github.com/microlib/simple"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	// create a key value map (to fake redis)
	store      map[string]string
	logger     simple.Logger
	connectors Clients
)

type Clients interface {
	LoginData(body []byte) (string, error)
	AllData(b []byte) ([]byte, error)
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
	Close() error
}

type FakeRedis struct {
}

type Connectors struct {
	// add mongodb connector here
	// mongodb *mgo
	Http  *http.Client
	redis FakeRedis
	name  string
}

type ConnectionData struct {
	Name          string
	RedisHost     string
	RedisPort     string
	HttpUrl       string
	MongoHost     string
	MongoPort     string
	MongoDatabase string
	MongoUser     string
	MongoPassword string
}

func (r *Connectors) Get(key string) (string, error) {
	return store[key], nil
}

func (r *Connectors) Set(key string, value string, expr time.Duration) (string, error) {
	store[key] = value
	return string(expr), nil
}

func (r *Connectors) Close() error {
	store = nil
	return nil
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestClients(data string, code int) Clients {

	// set all env variables for teting
	os.Setenv("VERSION", "1.1.0")
	os.Setenv("LOG_LEVEL", "info")

	logger.Level = os.Getenv("LOG_LEVEL")

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	file, _ := ioutil.ReadFile(data)
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(string(file))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	redisclient := FakeRedis{}
	conns := &Connectors{redis: redisclient, Http: httpclient, name: "test"}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestLoginData(t *testing.T) {

	var err error
	store = map[string]string{"hash": ""}
	store = map[string]string{"test": ""}

	// create anonymous struct
	tests := []struct {
		Name     string
		Payload  string
		Handler  string
		FileName string
		Want     int
		errorMsg string
	}{
		{
			"[TEST] Login with valid param and payload should pass",
			"{\"username\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMxNyA5QVQ=\"}",
			"LoginData",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] Login with invalid param and payload should fail",
			"{\"user\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"pass\":\"TFMxNyA5QVQ=\"}",
			"LoginData",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] Login with invalid payload should fail",
			"{\"",
			"LoginData",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] AllData with invalid response data should fail",
			"{\"username\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMxNyA5QVQ=\"}",
			"AllData",
			"tests/payload-errors.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] Alldata with valid param and payload should pass",
			"{\"apitoken\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMxNyA5QVQ=\"}",
			"AllData",
			"tests/payload-example.json",
			http.StatusOK,
			"Handler returned wrong status code got %d want %d",
		},
		{
			"[TEST] Alldata with invalid param and payload should fail",
			"{\"apit\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\"}",
			"AllData",
			"tests/payload-example.json",
			http.StatusInternalServerError,
			"Handler returned wrong status code got %d want %d",
		},
	}

	for _, tt := range tests {
		logger.Info(fmt.Sprintf("%s : \n", tt.Name))
		connectors = NewTestClients(tt.FileName, tt.Want)
		switch tt.Handler {
		case "LoginData":
			_, err = connectors.LoginData([]byte(tt.Payload))
		case "AllData":
			_, err = connectors.AllData([]byte(tt.Payload))
		}
		if err != nil && tt.Want == 200 {
			t.Errorf(tt.errorMsg, 200, 500)
		}
		if err == nil && tt.Want == 500 {
			t.Errorf(tt.errorMsg, 500, 200)
		}
	}
}
