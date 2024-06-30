FROM golang:1.22-alpine as build-env

ENV USER=user
ENV UID=1000

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/NA" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY go.mod go.sum .
COPY cmd/auth/main.go cmd/auth/
COPY pkg/ pkg/

RUN GOOS=linux GOARCH=amd64 go build -o bin/auth cmd/auth/main.go

FROM scratch

COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group
COPY --from=build-env /go/bin/auth /go/bin/auth

ENTRYPOINT ["/go/bin/auth"]