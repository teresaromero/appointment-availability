FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o checker .

FROM golang:1.24.0-alpine

RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates

WORKDIR /app
COPY --from=builder /app/checker .
CMD ["./checker"]