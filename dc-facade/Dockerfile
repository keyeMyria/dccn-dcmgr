FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/Ankr-network/dccn-hub/app_dccn_dcmgr
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i  -o app_dccn_dcmgr ./main.go

FROM scratch

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/app_dccn_dcmgr /
CMD ["/app_dccn_dcmgr"]
