
FROM golang:1.24 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o location_user_api ./cmd

FROM scratch

COPY --from=builder /app/location_user_api /location_user_api

ENTRYPOINT ["/location_user_api"]