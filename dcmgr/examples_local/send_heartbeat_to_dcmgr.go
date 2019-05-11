package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func main() {
    status := common_proto.DataCenterStatus{}
    status.DcStatus = common_proto.DCStatus_AVAILABLE



    status.DcId = "bb339361-e698-4d39-bf33-bbb2c7a8c308"
    status.DcHeartbeatReport = &common_proto.DCHeartbeatReport{}
    status.DcHeartbeatReport.Report = "this is report"
    status.DcHeartbeatReport.Metrics = "CPU 200 memory 3000"

    status.GeoLocation = &common_proto.GeoLocation{Lng:"9999", Lat:"8888"}



	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_HEARTBEAT,
		OpPayload: &common_proto.DCStream_DataCenter{DataCenter: &status},
	}

	log.Printf("messg %+v \n", &event)

	// publisher := micro2.NewPublisher("FromDcMgrToDcFacade")

	publisher := micro2.NewPublisher("FromDCFacadeToDCMgr") // appmgr to dcmgr

	publisher.Publish(&event)
}
