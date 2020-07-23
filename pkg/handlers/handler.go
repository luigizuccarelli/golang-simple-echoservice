package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-simple-oc4service/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-simple-oc4service/pkg/schema"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func EchoHandler(w http.ResponseWriter, r *http.Request, conn connectors.Clients) {
	var response *schema.Response
	var req *schema.Request

	// ensure we don't have nil - it will cause a null pointer exception
	if r.Body == nil {
		r.Body = ioutil.NopCloser(bytes.NewBufferString(""))
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "EchoHandler body data error %v"
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, err)
		fmt.Fprintf(w, string(b))
		return
	}

	conn.Trace("Request body : %s", string(body))

	// unmarshal result from mw backend
	errs := json.Unmarshal(body, &req)
	if errs != nil {
		msg := "EchoHandler could not unmarshal input data from servisBOT to schema %v"
		conn.Error(msg, errs)
		b := responseErrorFormat(http.StatusInternalServerError, w, msg, errs)
		fmt.Fprintf(w, string(b))
		return
	}
	response = &schema.Response{Code: http.StatusOK, Status: "OK", Message: req.Message}
	b, _ := json.MarshalIndent(response, "", "	")
	conn.Debug(fmt.Sprintf("EchoHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	addHeaders(w, r)
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \"golang-simple-oc4service\" }")
	return
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

// responsErrorFormat - utility function
func responseErrorFormat(code int, w http.ResponseWriter, msg string, val ...interface{}) []byte {
	var b []byte
	response := &schema.Response{Code: code, Status: "ERROR", Message: fmt.Sprintf(msg, val...)}
	w.WriteHeader(code)
	b, _ = json.MarshalIndent(response, "", "	")
	return b
}
