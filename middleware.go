package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	USERNAME     = "username"
	PASSWORD     = "password"
	APITOKEN     = "apitoken"
	ERRMSGFORMAT = " %s "
)

func (c Connectors) LoginData(body []byte) (string, error) {

	logger.Trace("In function LoginData")

	var j map[string]interface{}
	var schema SchemaInterface
	var apitoken = ""

	e := json.Unmarshal(body, &j)
	if e != nil {
		logger.Error(fmt.Sprintf("LoginData %v\n", e.Error()))
		return apitoken, e
	}
	if j[USERNAME] == nil || j[PASSWORD] == nil {
		return apitoken, errors.New("username and/or password field is nil")
	}
	logger.Debug(fmt.Sprintf("json data %v\n", j))

	// use this for cors
	//res.setHeader("Access-Control-Allow-Origin", "*")
	//res.setHeader("Access-Control-Allow-Methods", "POST")
	//res.setHeader("Access-Control-Allow-Headers", "accept, content-type")

	// lets first check in the in-memory cache
	key := j[USERNAME].(string) + ":" + j[PASSWORD].(string)
	//h := sha256.New()
	hashkey := sha256.Sum256([]byte(key))
	val, err := c.Get("hash")
	var newval [32]byte
	copy(newval[:], []byte(val))
	logger.Debug(fmt.Sprintf("Keys %x : %x", string(hashkey[:32]), val))
	if val == "" || err != nil {
		logger.Info(fmt.Sprintf("Key not found in cache %s", key))

		req, err := http.NewRequest("GET", config.Url+"/username/"+j[USERNAME].(string)+"/password/"+j[PASSWORD].(string), nil)
		req.Header.Set("token", config.Token)
		resp, err := c.Http.Do(req)
		logger.Info(fmt.Sprintf("Connected to host %s", config.Url))
		if err != nil || resp.StatusCode != 200 {
			logger.Error(fmt.Sprintf(ERRMSGFORMAT, err.Error()))
			return apitoken, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(fmt.Sprintf(ERRMSGFORMAT, err.Error()))
			return apitoken, err
		}
		logger.Trace(fmt.Sprintf("Response from MW %s", string(body)))

		errs := json.Unmarshal(body, &schema)
		if errs != nil {
			logger.Error(fmt.Sprintf(ERRMSGFORMAT, errs.Error()))
			return apitoken, errs
		}
		logger.Debug(fmt.Sprintf("Schema Data %v ", schema))

		//all, _ := json.MarshalIndent(schema, "", "	")
		_, err = c.Set("all", string(body), time.Hour)
		_, err = c.Set("hash", string(hashkey[:32]), time.Hour)
		_, err = c.Set(schema.Accounts[0].CustomerNumber, string(body), time.Hour)
		_, err = c.Set(schema.PostalAddresses[0].EmailAddress.EmailAddress, string(body), time.Hour)
		logger.Info(fmt.Sprintf("CustomerNumber %s", schema.Accounts[0].CustomerNumber))

		if err != nil {
			return apitoken, err
		}

		apitoken = xid.New().String()
		_, err = c.Set(APITOKEN, apitoken, time.Hour)
	} else if newval != hashkey {
		logger.Error(fmt.Sprintf("Hash token's don't match %s != %s", val, key))
		return apitoken, errors.New("hash token does not match")
	} else if hashkey == newval {
		logger.Info(fmt.Sprintf("Key found in cache %s", key))
		apitoken = xid.New().String()
		_, err = c.Set(APITOKEN, apitoken, time.Hour)
	}
	return apitoken, nil
}

func (c Connectors) AllData(b []byte) ([]byte, error) {

	logger.Trace("In function AllData")

	var subs []byte
	var data string
	var j map[string]interface{}

	e := json.Unmarshal(b, &j)
	if e != nil {
		logger.Error(fmt.Sprintf("AllData %v\n", e.Error()))
		return subs, e
	}
	logger.Debug(fmt.Sprintf("json data %v\n", j))

	// lets first check in the in-memory cache
	if j[APITOKEN] == nil {
		return subs, errors.New("Invalid or empty api token")
	}
	apitoken := j[APITOKEN].(string)
	val, err := c.Get(APITOKEN)
	logger.Info(fmt.Sprintf("Apitoken from cache %s : from req object %s", val, apitoken))
	if apitoken != val || err != nil {
		return subs, err
	} else {
		data, e = c.Get("all")
		if e != nil {
			return subs, e
		}
	}
	subs = []byte(data)
	return subs, nil
}

func (c Connectors) AllDataByCustomerNumber(customernumber string) ([]byte, error) {

	logger.Trace("In function AllDataByCustomerNumber")

	var subs []byte
	var data string

	// lets first check in the in-memory cache
	val, err := c.Get(customernumber)
	if val == "" {
		return subs, errors.New(fmt.Sprintf("CustomerNumber %s not found ", customernumber))
	} else {
		data, err = c.Get("all")
		if err != nil {
			return subs, err
		}
	}
	subs = []byte(data)
	return subs, nil
}
