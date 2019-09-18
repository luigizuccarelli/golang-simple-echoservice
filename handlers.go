package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"os"
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

	// lets first check in the in-memory cache
	key := j[USERNAME].(string) + ":" + j[PASSWORD].(string)
	hashkey := sha256.Sum256([]byte(key))
	val, err := c.Get(string(hashkey[:32]))
	logger.Debug(fmt.Sprintf("Key %x ", string(hashkey[:32])))
	if val == "" || err != nil {
		logger.Info(fmt.Sprintf("Key not found in cache %s", key))
		req, err := http.NewRequest("GET", os.Getenv("URL")+"/username/"+j[USERNAME].(string)+"/password/"+j[PASSWORD].(string), nil)
		req.Header.Set("token", os.Getenv("TOKEN"))
		resp, err := c.Http.Do(req)
		logger.Info(fmt.Sprintf("Connected to host %s", os.Getenv("URL")))
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

		errs := json.Unmarshal(body, &schema)
		if errs != nil {
			logger.Error(fmt.Sprintf(ERRMSGFORMAT, errs.Error()))
			return apitoken, errs
		}

		logger.Debug(fmt.Sprintf("Response from MW call %s", string(body)))
		_, err = c.Set(string(hashkey[:32]), string(body), time.Hour)
		_, err = c.Set("test", string(body), time.Hour)

		if err != nil {
			logger.Error(fmt.Sprintf(ERRMSGFORMAT, err.Error()))
			return apitoken, err
		}

		apitoken = xid.New().String()
		_, err = c.Set(APITOKEN, apitoken, time.Hour)
	} else {
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
	logger.Debug(fmt.Sprintf("Function AllData json data %v\n", j))

	// lets first check in the in-memory cache
	if j[APITOKEN] == nil {
		return subs, errors.New("Invalid or empty api token")
	}
	apitoken := j[APITOKEN].(string)
	val, err := c.Get(APITOKEN)
	logger.Trace(fmt.Sprintf("Function AllData : apitoken from cache %s : from req object %s", val, apitoken))
	if err != nil {
		logger.Error(fmt.Sprintf("Function AllData : error object %s ", err.Error()))
		return subs, err
	} else {
		data, e = c.Get("test")
		logger.Debug(fmt.Sprintf("Function AllData data from cache %s", data))
		if e != nil {
			logger.Error(fmt.Sprintf("Function AllData : error object %s ", e.Error()))
			return subs, e
		}
	}
	subs = []byte(data)
	return subs, nil
}
