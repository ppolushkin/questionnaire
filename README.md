# Questionnaire project

## Dependencies
dep ensure -add github.com/pkg/errors 

## Docker
https://medium.com/travis-on-docker/how-to-dockerize-your-go-golang-app-542af15c27a2 

## Build for docker
CGO_ENABLED=0 go build -a -installsuffix cgo

## PROD Run all
go clean && CGO_ENABLED=0 go build -a -installsuffix cgo && sudo docker-compose -f docker-compose-prod.yml up  

## DEV Run db 
go clean && sudo docker-compose -f docker-compose-dev.yml up

then right button of folder and Debug

## Run integration tests
sudo docker-compose -f docker-compose-test.yml up

then right click on main_test.go

#Before commit

## Format code
gofmt -w .

## Check code 
golint .

### Questionnaire application
http://localhost:8080

### PHP my admin
http://localhost:8001


# TODO:

- implement logic

- logging (kubernates?)
- monitoring (kubernates?)
- swagger
- flyway
- angularjs on nodejs or on go
- deploy to cloud


- aws? swiss com? DO? asure?
- kubernates


# Links 
https://medium.com/@kelvin_sp/building-and-testing-a-rest-api-in-golang-using-gorilla-mux-and-mysql-1f0518818ff6 

https://docs.docker.com/get-started/

