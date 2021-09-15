test:
	go test ./... -p=8 -count=4 -race

run:
	go run cmd/feeder/main.go

build:
	go build -o feeder cmd/feeder/main.go

run-client:
	go run cmd/client/main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
