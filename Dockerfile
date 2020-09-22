FROM golang:1.13.6-alpine AS builder

ARG SERVICE

RUN apk add --update --no-cache git build-base musl-dev linux-headers
RUN mkdir /build
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o bin/watchmarket ./cmd/$SERVICE

FROM alpine:latest
COPY --from=builder /build/bin /app/bin/$SERVICE
COPY --from=builder /build/config.yml /config/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/bin/watchmarket", "-c", "/config/config.yml"]
