package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"lmzsoftware.com/lzuccarelli/golang-simple-echoservice/pkg/connectors"
	"lmzsoftware.com/lzuccarelli/golang-simple-echoservice/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func EchoHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var req *schema.Request

	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "EchoHandler body data error %v"
		b := responseFormat(http.StatusInternalServerError, "KO", w, msg, err)
		fmt.Fprintf(w, "%s", b)
		return
	}

	conn.Trace("Request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &req)
	if errs != nil {
		msg := "EchoHandler could not unmarshal input data from json to schema %v"
		conn.Error(msg, errs)
		b := responseFormat(http.StatusInternalServerError, "KO", w, msg, errs)
		fmt.Fprintf(w, "%s", b)
		return
	}
	// simulate an error - used for metrics
	if req.Id == "error" {
		msg := "EchoHandler simulate 500 internal server error %s"
		conn.Error(msg, req.Message)
		b := responseFormat(http.StatusInternalServerError, "KO", w, msg, "forced error")
		fmt.Fprintf(w, "%s", b)
		return
	}
	response := responseFormat(http.StatusOK, "OK", w, req.Message)
	conn.Debug("EchoHandler response : %s", response)
	fmt.Fprintf(w, "%s", response)
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("API-KEY") != "" {
		w.Header().Set("API_KEY_PT", r.Header.Get("API_KEY"))
	}
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// responsFormat - utility function
func responseFormat(code int, status string, w http.ResponseWriter, msg string, val ...interface{}) string {
	response := `{"Code":"` + strconv.Itoa(code) + `", "Status": "` + status + `", "Message":"` + fmt.Sprintf(msg, val...) + `"}`
	w.WriteHeader(code)
	return response
}
