package main

import (
	"crypto/tls"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

type Clients interface {
	LoginData(body []byte) (string, error)
	AllData(b []byte) ([]byte, error)
	AllDataByCustomerNumber(customernumber string) ([]byte, error)
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
	Close() error
}

type Connectors struct {
	// add mongodb connector here
	// mongodb *mgo
	Http  *http.Client
	redis *redis.Client
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

func NewClientConnectors(conn ConnectionData) Clients {
	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:         conn.RedisHost + ":" + conn.RedisPort,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Password:     "",
		DB:           0,
	})

	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	conns := &Connectors{redis: redisClient, Http: httpClient, name: conn.Name}
	return conns
}

func (r *Connectors) Get(key string) (string, error) {
	val, err := r.redis.Get(key).Result()
	return val, err
}

func (r *Connectors) Set(key string, value string, expr time.Duration) (string, error) {
	val, err := r.redis.Set(key, value, expr).Result()
	return val, err
}

func (r *Connectors) Close() error {
	r.redis.Close()
	return nil
}
