FROM golang:1.23.1-alpine AS builder

RUN /sbin/apk update && \
    /sbin/apk --no-cache add ca-certificates git tzdata && \
    /usr/sbin/update-ca-certificates

RUN  apk add --no-cache build-base

RUN adduser -D -g '' products
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env ./
COPY cmd   ./cmd
COPY pkg ./pkg


RUN CGO_ENABLED=1 go build -a -tags netgo,osusergo,musl \
    -ldflags "-extldflags '-static' -s -w" \
    -o products-micro ./cmd

FROM busybox:musl
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/products-micro /app
COPY --from=builder /app/.env /app

USER products

EXPOSE 8000
ENTRYPOINT ["./products-micro"]