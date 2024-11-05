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

test-docker: 
	@docker container stop testing.local.db # Stop container if running
	docker compose -f test.docker-compose.yml up -d database
	go run ./cmd/migrate -env test.env up
	docker compose -f test.docker-compose.yml up testing
	go run ./cmd/migrate -env test.env reset
	docker container stop testing.local.db

format:
 	$('go fmt ./pkg/...')
