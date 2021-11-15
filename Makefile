

.PHONY: swagger
swagger:
	GO111MODULE=off swagger generate spec -o ./app/docs/swagger.yaml --scan-models
	swagger serve -F=swagger ./app/docs/swagger.yaml

.PHONY: run
run:
	go mod tidy
	go run main.go
