.PHONY: run build test migrate-up migrate-down

run:
	go run main.go

build:
	go build -o bin/kaafipay-backend main.go

test:
	go test -v ./...

migrate-up:
	migrate -path internal/db/migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path internal/db/migrations -database "$${DATABASE_URL}" down

mock:
	mockgen -source=internal/repository/user_repository.go -destination=internal/mocks/user_repository_mock.go 