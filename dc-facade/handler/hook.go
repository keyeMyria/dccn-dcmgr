package handler

import (
	"log"
	"net"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/hashicorp/yamux"
)

func (p *Relay) OnAccept(key string, conn net.Conn) error      { return nil }
func (p *Relay) OnBuild(key string, conn *yamux.Session) error { return nil }

func (p *Relay) OnClose(key string, conn *yamux.Session) {
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
