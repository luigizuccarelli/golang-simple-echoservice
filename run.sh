#!/bin/bash

LOG_LEVEL=trace
SERVER_PORT=9000
VERSION=1.0.1
NAME=golang-simple-oc4service

export LOG_LEVEL NAME SERVER_PORT VERSION 

./build/microservice
