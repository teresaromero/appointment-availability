FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o checker .

FROM scratch

WORKDIR /app
COPY --from=builder /app/checker .
CMD ["./checker"]