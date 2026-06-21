FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s' \
    -a \
    -o /build/bin/topsigner \
    ./cmd/topsigner

FROM alpine:3.21

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /build/bin/topsigner /usr/local/bin

USER nonroot

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/topsigner" ]
