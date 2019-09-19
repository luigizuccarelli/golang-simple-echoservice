# MyPorfolio middleware interface golang microservice

A simple golang microservice that interfaces with the middleware api.

## Status 
[![Quality gate](https://sonarqube-cicd.apps.aws2-dev.ocp.14west.io/api/project_badges/quality_gate?project=myportfolio-middlewareinterface)](https://sonarqube-cicd.apps.aws2-dev.ocp.14west.io/dashboard?id=myportfolio-middlewareinterface)

## Coverage Info

To view the current coverage info please click on the [link](https://sonarqube-cicd.apps.aws2-dev.ocp.14west.io/dashboard?id=myportfolio-middlewareinterface)


## Usage 

```bash
# cd to project directory and build executable
$ go build -o bin/microservice .

```

## Docker build

```bash
docker build -t <your-registry-id>/myportfolio-middlewareinterface:1.11.0 .

```

## Curl timing usage
```
curl -w "@curl-timing.txt" -o /dev/null -s "http://site-to-test

```

## Executing tests
```bash
# clear the cache - this is optional
go clean -testcache
GOCACHE=off go test -v config.go config_test.go schema.go handlers.go middleware.go middleware_test.go handlers_test.go -coverprofile tests/results/cover.out
go tool cover -html=tests/results/cover.out -o tests/results/cover.html
# run sonarqube scanner (assuming sonarqube server is running)
# NB the SonarQube host and login will differ - please update it accordingly 
 ~/Programs/sonar-scanner-3.3.0.1492-linux/bin/sonar-scanner  -Dsonar.projectKey=myportfolio-middlewareinterface  -Dsonar.sources=.   -Dsonar.host.url=http://<add-host-here>   -Dsonar.login=<add-login-token-here> -Dsonar.go.coverage.reportPaths=tests/results/cover.out -Dsonar.exclusions=vendor/**,*_test.go,main.go,connectors.go,tests/**

```
## Testing container 
```bash

# start the container
# curl the isalive endpoint
curl -k -H 'Token: xxxxx' -w '@curl-timing.txt'  http://127.0.0.1:9000/api/v1/sys/info/isalive

```
