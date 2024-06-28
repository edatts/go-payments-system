build:
	go build -o bin/auth cmd/auth/main.go

test:
	go test -cover -coverprofile=c.out ./...

run: build
	./bin/auth

run-db-dev:
	sudo docker run -it --rm -p 5432:5432 -d -e POSTGRES_PASSWORD=postgres --name postgres postgres

stop-db-dev:
	sudo docker stop postgres

install-migrate: # Replace postgres tag to build for other db
	go install -tags 'postgres'  github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

migration:
	migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down