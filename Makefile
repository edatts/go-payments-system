build:
	go build -o bin/auth cmd/auth/main.go

test:
	go test -v -cover -coverprofile=c.out ./...

run: build
	./bin/auth