FROM golang:1.11.4-alpine as builder

RUN apk add -U --no-cache ca-certificates

WORKDIR /go/src/github.com/Ankr-network/dccn-hub

COPY . .

ARG APP_DOMAIN

RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -i -o cmd/email \
    -ldflags "-X github.com/Ankr-network/dccn-hub/app-dccn-email/subscriber.APPDOMAIN=${APP_DOMAIN}" \
    app-dccn-email/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/cmd/email /

CMD ["/email"]
