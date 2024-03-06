SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build the docker file"
	@echo "	   build-multiplatform			Build the docker file for multiple platforms"
	@echo "	   run					Launch a working cache and voter api"
	@echo "	   stop					Stop a working cache and voter api"
	@echo "	   test					Run the project tests against a stood up api"
	@echo "	   get-all				Get all voters"
	@echo "	   delete-all				Delete all voters"

.PHONY: build
build: Dockerfile
	docker build -t voter-api .

.PHONY: build-multiplatform
build-multiplatform: Dockerfile
	docker buildx build --platform linux/amd64,linux/arm64 -t voter-api .


.PHONY: run
run: build
	docker compose up -d

stop:
	docker compose down

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: get-all
get-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters

.PHONY: delete-all
delete-all:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters