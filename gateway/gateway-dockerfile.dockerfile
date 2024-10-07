FROM golang:1.23.1-alpine AS builder

RUN /sbin/apk update && \
    /sbin/apk --no-cache add ca-certificates git tzdata && \
    /usr/sbin/update-ca-certificates

RUN  apk add --no-cache build-base
RUN adduser -D -g '' gateway
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./


RUN CGO_ENABLED=0 go build \
    -ldflags "-extldflags '-static' -s -w" \
    -o gateway ./

FROM busybox:musl
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/gateway /app
COPY --from=builder /app/.env /app

USER gateway

EXPOSE 80
ENTRYPOINT ["./gateway"]