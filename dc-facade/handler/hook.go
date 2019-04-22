package handler

import (
	"log"
	"net"

	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func (p *Relay) OnAccept(key string, conn net.Conn) error      { return nil }
func (p *Relay) OnBuild(key string, conn *pgrpc.Session) error { return nil }

func (p *Relay) OnClose(key string, conn *pgrpc.Session) {
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
