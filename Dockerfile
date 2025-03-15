FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o checker cmd/checker/main.go

FROM alpine:3.21.3

RUN apk add --no-cache \
    chromium \
    harfbuzz \
    nss \
    freetype \
    ttf-freefont \
    && rm -rf /var/cache/apk/*

# Set Chrome as default
ENV CHROME_BIN=/usr/bin/chromium-browser

WORKDIR /app
COPY --from=builder /app/checker .
CMD ["./checker"]