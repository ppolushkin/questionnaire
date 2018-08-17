# Questionnaire project

## Dependencies
dep ensure -add github.com/pkg/errors 

## Docker
https://medium.com/travis-on-docker/how-to-dockerize-your-go-golang-app-542af15c27a2 

## Build for docker
CGO_ENABLED=0 go build -a -installsuffix cgo

# PROD Run all
in main.go host should be db

go clean && CGO_ENABLED=0 go build -a -installsuffix cgo && sudo docker-compose -f docker-compose-prod.yml up  


# DEV Run db
in main.go host should be localhost 

go clean && sudo docker-compose -f docker-compose-dev.yml up

then right button of folder and Debug

# Run integration tests

sudo docker-compose -f docker-compose-test.yml up

then right click on main_test.go

### Questionnaire application
http://localhost:8080

### PHP my admin
http://localhost:8001









# Not used

## Build "quest" docker image
sudo docker build -t quest .

## Run docker
sudo docker run -p 8080:8080 quest

## Notes about build

env GOOS=linux GARCH=amd64 go install -v -a -tags netgo -installsuffix netgo -ldflags "-linkmode external -extldflags -static"

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo 

go build --ldflags '-w -linkmode external -extldflags "-static"' main.go
