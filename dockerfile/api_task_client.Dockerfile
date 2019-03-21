FROM golang:1.11.4-alpine as builder

WORKDIR /go/src/github.com/Ankr-network/dccn-hub
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd/api_taskmgr_client app-dccn-api/examples/task.go

FROM scratch

COPY --from=builder /go/src/github.com/Ankr-network/dccn-hub/cmd/api_taskmgr_client  /
CMD ["/api_taskmgr_client"]
