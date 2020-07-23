package main

import (
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-simple-oc4service/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-simple-oc4service/pkg/handlers"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/golang-simple-oc4service/pkg/validator"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

var (
	logger *simple.Logger
)

func startHttpServer(con connectors.Clients) *http.Server {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/echo", func(w http.ResponseWriter, req *http.Request) {
		handlers.EchoHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/sys/info/isalive", handlers.IsAlive).Methods("GET")

	sh := http.StripPrefix("/api/v1/api-docs", http.FileServer(http.Dir("./swaggerui")))
	r.PathPrefix("/api/v1/api-docs").Handler(sh)

	http.Handle("/", r)

	if err := srv.ListenAndServe(); err != nil {
		con.Error("Httpserver: ListenAndServe() error: " + err.Error())
	}

	return srv
}

func main() {
	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}
	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}
	conn := connectors.NewClientConnections(logger)
	srv := startHttpServer(conn)
	logger.Info("Starting server on port " + srv.Addr)
}
