FROM golang:1.23.1-alpine AS builder

RUN /sbin/apk update && \
    /sbin/apk --no-cache add ca-certificates git tzdata && \
    /usr/sbin/update-ca-certificates

RUN  apk add --no-cache build-base

RUN adduser -D -g '' orders
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY .env ./
COPY cmd   ./cmd
COPY pkg ./pkg


RUN CGO_ENABLED=1 go build -tags musl \
    -ldflags "-extldflags '-static' -s -w" \
    -o orders-micro ./cmd

FROM busybox:musl
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/orders-micro /app
COPY --from=builder /app/.env /app

USER orders

EXPOSE 4000
ENTRYPOINT ["./orders-micro"]