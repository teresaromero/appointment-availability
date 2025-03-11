FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o checker .

FROM scratch

WORKDIR /app
COPY --from=builder /app/checker .
CMD ["./checker"]