FROM golang:1.22-alpine AS build-env

COPY cmd/migrate/main.go cmd/migrate/main.go
COPY cmd/migrate/migrations cmd/migrate/migrations
COPY pkg/config pkg/config
COPY go.mod go.sum .

RUN CGO_ENABLED=0 go build -o bin/migrate cmd/migrate/main.go

# COPY bin/migrate bin/migrate

FROM scratch

WORKDIR /go

COPY --from="build-env" /go/cmd/migrate/migrations/ /go/cmd/migrate/migrations/
COPY --from="build-env" /go/bin/migrate /go/bin/migrate

ENTRYPOINT ["/go/bin/migrate", "up"]