package handler

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

type PgrpcHook struct {
	taskFeedback *micro2.Publisher
	*pgrpc.ProxyProtoHook
}

func NewPgrpcHook(feedback *micro2.Publisher) *PgrpcHook {
	handler := &PgrpcHook{
		taskFeedback:   feedback,
		ProxyProtoHook: &pgrpc.ProxyProtoHook{},
	}
	return handler
}

func (p *PgrpcHook) OnClose(key string, conn *pgrpc.Session) {
	log.Printf("public %s close message", key)

	p.taskFeedback.Publish(&common_proto.DCStream{
		OpType: common_proto.DCOperation_HEARTBEAT,
		OpPayload: &common_proto.DCStream_DataCenter{
			DataCenter: &common_proto.DataCenterStatus{
				DcId:     key,
				DcStatus: common_proto.DCStatus_UNAVAILABLE,
			},
		},
	})
}
