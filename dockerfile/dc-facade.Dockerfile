FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/Ankr-network/dccn-dcmgr
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i  -o cmd/dc-facade dc-facade/main.go

ADD https://s3-us-west-1.amazonaws.com/static.ankr.com/geo/GeoIP2-City.mmdb .

CMD ["cmd/dc-facade"]
