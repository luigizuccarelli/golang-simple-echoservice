package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	logger     simple.Logger
	config     Config
	connectors Clients
)

func startHttpServer(cfg Config) *http.Server {

	config = cfg

	logger.Debug(fmt.Sprintf("Config in startServer %v ", cfg))
	srv := &http.Server{Addr: ":" + cfg.Port}

	r := mux.NewRouter()
	r.HandleFunc("/v1/sys/info/isalive", IsAlive).Methods("GET")
	r.HandleFunc("/v1/login", MiddlewareLogin).Methods("POST")
	r.HandleFunc("/v1/alldata", MiddlewareData).Methods("POST")
	r.HandleFunc("/v2/postaladdress/customernumber/{customerNumber}", MiddlewareCustomerNumberData).Methods("GET")
	http.Handle("/", r)

	connectionData := ConnectionData{
		Name:          "RealConnector",
		RedisHost:     cfg.RedisDB.Host,
		RedisPort:     cfg.RedisDB.Port,
		HttpUrl:       cfg.Url,
		MongoHost:     cfg.MongoDB.Host,
		MongoPort:     cfg.MongoDB.Port,
		MongoDatabase: cfg.MongoDB.DatabaseName,
		MongoUser:     cfg.MongoDB.User,
		MongoPassword: cfg.MongoDB.Password,
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
	// read the config
	config, _ := Init("config.json")
	logger.Level = config.Level
	srv := startHttpServer(config)
	logger.Info("Starting server on port " + config.Port)
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
