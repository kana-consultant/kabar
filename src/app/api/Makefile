.PHONY: run build migrate-up migrate-down clean

run:
	go run cmd/api/main.go

build:
	go build -o bin/seo-api cmd/api/main.go

migrate-up:
	go run cmd/api/main.go -migrate=up

migrate-down:
	go run cmd/api/main.go -migrate=down

clean:
	rm -rf bin/