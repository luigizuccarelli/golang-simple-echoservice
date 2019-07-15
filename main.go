package main

import (
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger     simple.Logger
	connectors Clients
)

func startHttpServer() *http.Server {

	srv := &http.Server{Addr: ":" + os.Getenv("PORT")}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sys/info/isalive", IsAlive).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/login", MiddlewareLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/alldata", MiddlewareData).Methods("POST", "OPTIONS")
	http.Handle("/", r)

	connectionData := ConnectionData{
		Name:          "RealConnector",
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		HttpUrl:       os.Getenv("URL"),
		MongoHost:     os.Getenv("MONGO_HOST"),
		MongoPort:     os.Getenv("MONGO_PORT"),
		MongoDatabase: os.Getenv("MONGO_DATABASE"),
		MongoUser:     os.Getenv("MONGO_USER"),
		MongoPassword: os.Getenv("MONGO_PWD"),
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
