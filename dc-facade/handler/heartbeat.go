package handler

import (
	"context"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	sign "github.com/Ankr-network/dccn-common/cert/sign"
	"github.com/Ankr-network/dccn-common/pgrpc"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"google.golang.org/grpc"
	"log"
	"time"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
)

type HeartBeat struct {
	taskFeedback *micro2.Publisher
	db dbservice.DBService
}

func NewHeartBeat(feedback *micro2.Publisher) *HeartBeat {
	dbInstance, _ := dbservice.New()
	handler := &HeartBeat{
		taskFeedback: feedback,
		db: dbInstance,
	}
	return handler
}


func (p *HeartBeat) StartCollectStatus() {
	for range time.Tick(20 * time.Second) {
        p.CheckEachDatacenterConnections()
		log.Printf("heartbeat finish")
	}
}


func (p *HeartBeat) CheckEachDatacenterConnections(){
	pgrpc.Each(func(key string, conn *grpc.ClientConn, err error) {
		// handle dial error
		status, ok := p.CheckDatacenterConnectionOK(key, conn, err)

		if ok {
			if key != status.DcId {
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

func (p *HeartBeat) CheckDatacenterConnectionOK(key string, conn *grpc.ClientConn, err error) (*common_proto.DataCenterStatus,  bool) {
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
	}else{
		log.Printf("not pass RsaVerify \n")
	}



	// FIXME: transaction
	// update status into db
	if rsp.ClusterId == "" {
		log.Printf("error for datacenter id does not exist")
		return nil, false
	}

	return rsp.Status, true
}
