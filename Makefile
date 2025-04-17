prepare:
	curl -s -o /dev/null https://pre-commit.com/install-local.py
	pre-commit install

generate: grpc
	rm -rf mocks/
	go generate tools.go
	go mod tidy

migrate-up:
	go run cmd/db/up/main.go

docker-up:
	docker compose -p integration-test -f docker-compose.integration.yaml up -d
	sleep 10
	make migrate-up

docker-down:
	docker compose -p integration-test -f docker-compose.integration.yaml down

integration_test:
	go test -tags integration_test -coverprofile=coverage.txt -covermode=atomic ./...

test_and_coverage:
	go test -race -covermode=atomic ./...

grpc:
	rm -rf generated && mkdir -p generated
	protoc --go_out=generated --go_opt=paths=source_relative --go-grpc_out=generated --go-grpc_opt=paths=source_relative proto/candle.proto

create-migration:
	go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir internal/migrations/scripts $(name)
