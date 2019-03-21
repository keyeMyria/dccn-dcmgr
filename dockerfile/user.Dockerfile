FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/Ankr-network/dccn-hub
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -i -o cmd/usermgr app-dccn-usermgr/main.go

FROM scratch

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/cmd/usermgr /
CMD ["/usermgr"]
