# user-management
This is a user management REST service.

## Installation

1. Run
   - `make swagger` to start swagger ui to test apis (needs swagger to be installed)
   - `go run main.go` to start microservice in on local machine

2. Deploy in `docker`
   - `docker-compose up -d` to build and start microservice in docker with postgress

3. Install swagger
   - `brew tap go-swagger/go-swagger`
   - `brew install go-swagger`
