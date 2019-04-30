package main

import (
	"bytes"
	"github.com/microlib/simple"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	// create a key value map (to fake redis)
	store      map[string]string
	logger     simple.Logger
	config     Config
	connectors Clients
)

type Clients interface {
	LoginData(body []byte) (string, error)
	AllData(b []byte) ([]byte, error)
	AllDataByCustomerNumber(customernumber string) ([]byte, error)
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

	// read the config
	config, _ := Init("config.json")
	logger.Level = config.Level

	// initialise our store (cache)
	store = map[string]string{"hash": ""}
	store = map[string]string{"all": ""}

	connectors = NewTestClients("tests/payload-example.json", 200)
	json := "{\"username\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMxNyA5QVQ=\"}"
	apitoken, err := connectors.LoginData([]byte(json))
	//fmt.Println(fmt.Sprintf("Body %v", body))
	if err != nil {
		t.Errorf("The login test should pass")
	}

	// test for failure
	json = "{\"{\"}"
	_, err = connectors.LoginData([]byte(json))
	if err == nil {
		t.Errorf("The login test should fail")
	}

	// test for failure
	json = "{\"user\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"pass\":\"TFMxNyA5QVQ=\"}"
	_, err = connectors.LoginData([]byte(json))
	if err == nil {
		t.Errorf("The login test should fail")
	}

	// test for failure
	connectors = NewTestClients("tests/config-errors.json", 200)
	json = "{\"username\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMQVQ=\"}"
	_, err = connectors.LoginData([]byte(json))
	if err == nil {
		t.Errorf("The login test should fail")
	}

	// test for pass
	connectors = NewTestClients("tests/payload-example.json", 200)
	json = "{\"username\": \"MTBLUlVOTkVSQFRBTEtUQUxLLk5FVA==\",\"password\":\"TFMxNyA5QVQ=\"}"
	apitoken, err = connectors.LoginData([]byte(json))
	if err != nil {
		t.Errorf("The login test should pass")
	}

	// test for pass
	json = "{\"apitoken\": \"dfsfdsfdsfsdf\"}"
	_, err = connectors.AllData([]byte(json))
	if err != nil {
		t.Errorf("The alldata test should pass")
	}

	// test for pass
	json = "{\"apitoken\": \"" + apitoken + "\"}"
	_, err = connectors.AllData([]byte(json))
	if err != nil {
		t.Errorf("The alldata test should pass")
	}

	// test for fail
	json = ""
	_, err = connectors.AllData([]byte(json))
	if err == nil {
		t.Errorf("The alldata test should fail")
	}

	// test for fail
	json = "{\"{\": \"}\"}"
	_, err = connectors.AllData([]byte(json))
	if err == nil {
		t.Errorf("The alldata test should fail")
	}

	_, err = connectors.AllDataByCustomerNumber("")
	if err == nil {
		t.Errorf("The alldatabycustomernumber test should fail")
	}

	// should pass
	_, err = connectors.AllDataByCustomerNumber("000002096234")
	if err != nil {
		t.Errorf("The alldatabycustomernumber test should pass")
	}

}
