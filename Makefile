build: 
	go build -o bin/api cmd/app/main.go

run: build
	bin/api

dev:
	air -c .air.toml

up:
	go run cmd/migrate/main.go up

down:
	go run cmd/migrate/main.go down

reset:
	go run cmd/migrate/main.go reset

test:
	go test ./pkg/...

test-verbose:
	go test ./pkg/... -v

format:
 	go fmt ./pkg/...
