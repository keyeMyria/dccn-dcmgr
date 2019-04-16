package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func main() {
    status := common_proto.DataCenterStatus{}
    status.Status = common_proto.DCStatus_AVAILABLE



    status.Id = "b0c6cba3-d407-41ec-b6e7-c5517d06c4a7"
    status.Name = "letian datacenter number 1"
    status.DcHeartbeatReport = &common_proto.DCHeartbeatReport{}
    status.DcHeartbeatReport.Report = "123"
    status.DcHeartbeatReport.Metrics = "11111"

	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_HEARTBEAT,
		OpPayload: &common_proto.DCStream_DataCenter{DataCenter: &status},
	}

	log.Printf("messg %+v \n", &event)

	// publisher := micro2.NewPublisher("FromDcMgrToDcFacade")

	publisher := micro2.NewPublisher("FromDCFacadeToDCMgr") // appmgr to dcmgr

	publisher.Publish(&event)
}
