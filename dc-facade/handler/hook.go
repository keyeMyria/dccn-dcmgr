package handler

import (
	"log"
	"net"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

type PGRPCCallBack struct {
	taskFeedback *micro2.Publisher
}

func NewPGRPCCallBack(feedback *micro2.Publisher) *PGRPCCallBack {
	handler := &PGRPCCallBack{
		taskFeedback: feedback,
	}
	return handler
}


func (p *PGRPCCallBack) OnAccept(key string, conn net.Conn) error      { return nil }
func (p *PGRPCCallBack) OnBuild(key string, conn *pgrpc.Session) error { return nil }

func (p *PGRPCCallBack) OnClose(key string, conn *pgrpc.Session) {
	log.Printf("public %s close message", key)

	p.taskFeedback.Publish(&common_proto.DCStream{
		OpType: common_proto.DCOperation_HEARTBEAT,
		OpPayload: &common_proto.DCStream_DataCenter{
			DataCenter: &common_proto.DataCenterStatus{
				Id:     key,
				Status: common_proto.DCStatus_UNAVAILABLE,
			},
		},
	})
}
