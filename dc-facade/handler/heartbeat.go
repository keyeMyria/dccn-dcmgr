package handler

import (
	"context"
	"log"
	"sync"
	"time"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/cert/sign"
	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	geo "github.com/Ankr-network/dccn-dcmgr/dc-facade/geo"
	dbservice "github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
	"google.golang.org/grpc"
)

type HeartBeat struct {
	taskFeedback *micro2.Publisher
	db           dbservice.DBService
}

func NewHeartBeat(feedback *micro2.Publisher) *HeartBeat {
	dbInstance, _ := dbservice.New()
	handler := &HeartBeat{
		taskFeedback: feedback,
		db:           dbInstance,
	}
	return handler
}

func (p *HeartBeat) StartCollectStatus() {
	var ipCache = &sync.Map{}
	for range time.Tick(20 * time.Second) {
		p.CheckEachDatacenterConnections(ipCache)
		log.Printf("heartbeat finish")
	}
}

func (p *HeartBeat) CheckEachDatacenterConnections(ipCache *sync.Map) {
	log.Printf("CheckEachDatacenterConnections  start \n")
	pgrpc.Each(func(key string, conn *grpc.ClientConn, err error) {
		// handle dial error
		ip := key
		if val, ok := ipCache.Load(key); ok {
			ip = val.(string)
		}
		log.Printf("key %s ip %s  start \n", key, ip)
		status, ok := p.CheckDatacenterConnectionOK(key, ip, conn, err)

		if ok {
			if key != status.DcId {
				ipCache.Store(status.DcId, ip)
				log.Printf("alias %s into %s", key, status.DcId)
				pgrpc.Alias(key, status.DcId, true)
			}

			event := common_proto.DCStream{
				OpType: common_proto.DCOperation_HEARTBEAT,
				OpPayload: &common_proto.DCStream_DataCenter{
					DataCenter: status,
				},
			}

			log.Println("publishing status of key ", key)
			p.taskFeedback.Publish(&event)
		}
	})
}

func (p *HeartBeat) CheckDatacenterConnectionOK(key, ip string, conn *grpc.ClientConn, err error) (*common_proto.DataCenterStatus, bool) {
	if err != nil {
		log.Println(err)
		return nil, false
	}
	defer conn.Close()

	// collect status(heartbeat)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("collecting status of key ", key)
	timestamp := time.Now()
	timestampStr := timestamp.String()
	micro2.WriteLog("timestampe " + timestampStr)
	rsp, err := dcmgr.NewDCClient(conn).Overview(ctx, &dcmgr.DCOverviewRequest{Timestamp: timestampStr})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	dcID := rsp.ClusterId
	status, _ := p.db.GetByID(dcID)

	client_cert := status.Clientcert

	if sign.RsaVerify(client_cert, timestampStr, rsp.Signature) {
		log.Printf("pass RsaVerify  timestampStr  %s %s %s  \n", client_cert, timestampStr, rsp.Signature)
	} else {
		log.Printf("not pass RsaVerify \n")
	}

	// FIXME: transaction
	// update status into db
	if rsp.ClusterId == "" {
		log.Printf("error for datacenter id does not exist")
		return nil, false
	}

	lat, lang, country := geo.ReadIPinfo(ip)
	log.Printf("%s %s %s ", lat, lang, country)
	location := common_proto.GeoLocation{}
	location.Lat = lat
	location.Lng = lang
	location.Country = country
	rsp.Status.GeoLocation = &location

	return rsp.Status, true
}
