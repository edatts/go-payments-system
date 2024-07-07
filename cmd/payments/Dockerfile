FROM golang:1.22-alpine AS build-env

ENV USER=nonroot
ENV UID=1000

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/NA" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /root

COPY go.mod go.sum .
RUN go mod download
# RUN go get -u ./...

COPY cmd/payments/main.go cmd/payments/
COPY pkg/ pkg/

RUN GOOS=linux GOARCH=amd64 go build -o bin/payments cmd/payments/main.go

RUN ls -al /root
RUN ls -al /root/bin

FROM scratch

COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group
COPY --from=build-env /root/bin/payments /bin/payments

USER "${USER}"

ENTRYPOINT ["/bin/payments"]