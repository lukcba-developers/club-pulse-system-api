.PHONY: run build test test-coverage docker-build docker-run

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

test-coverage:
	./scripts/coverage.sh

docker-build:
	docker build -t club-pulse-backend .

docker-run:
	docker run -p 8080:8080 club-pulse-backend

swag:
	swag init -g cmd/api/main.go -o docs
