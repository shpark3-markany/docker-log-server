# Build stae
FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .


# Run stage
FROM alpine:3.13

COPY --from=builder /app/main /main

CMD ["/main", "server"]