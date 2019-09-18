package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	logger     simple.Logger
	connectors Clients
)

func startHttpServer() *http.Server {

	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sys/info/isalive", IsAlive).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/login", MiddlewareLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/alldata", MiddlewareData).Methods("POST", "OPTIONS")
	http.Handle("/", r)

	connectionData := ConnectionData{
		Name:          "RealConnector",
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		HttpUrl:       os.Getenv("URL"),
		MongoHost:     "",
		MongoPort:     "",
		MongoDatabase: "",
		MongoUser:     "",
		MongoPassword: "",
	}

	connectors = NewClientConnectors(connectionData)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {
	ValidateEnvars()
	// read the log level
	logger.Level = os.Getenv("LOG_LEVEL")
	srv := startHttpServer()
	logger.Info("Starting server on port " + os.Getenv("PORT"))
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan

	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown successfully")
	os.Exit(code)
}

func checkEnvar(item string) {
	name := strings.Split(item, ",")[0]
	required, _ := strconv.ParseBool(strings.Split(item, ",")[1])
	if os.Getenv(name) == "" {
		if required {
			logger.Error(fmt.Sprintf("%s envar is mandatory please set it", name))
			os.Exit(-1)
		} else {
			logger.Error(fmt.Sprintf("%s envar is empty please set it", name))
		}
	}
}

// ValidateEnvars : public call that groups all envar validations
// These envars are set via the openshift template
func ValidateEnvars() {
	items := []string{
		"LOG_LEVEL,false",
		"SERVER_PORT,true",
		"REDIS_HOST,true",
		"REDIS_PORT,true",
		"REDIS_PASSWORD,true",
		"VERSION,true",
		"URL,true",
	}
	for x, _ := range items {
		checkEnvar(items[x])
	}
}
