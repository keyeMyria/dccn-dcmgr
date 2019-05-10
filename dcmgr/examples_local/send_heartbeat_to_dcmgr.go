package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func main() {
    status := common_proto.DataCenterStatus{}
    status.DcStatus = common_proto.DCStatus_AVAILABLE



    status.DcId = "5da3c264-6433-4ca4-897a-1d5e5c02cf16"
    status.DcHeartbeatReport = &common_proto.DCHeartbeatReport{}
    status.DcHeartbeatReport.Report = "this is report"
    status.DcHeartbeatReport.Metrics = "CPU 200 memory 3000"

    status.GeoLocation = &common_proto.GeoLocation{Lng:"lag xxxx", Lat:"22222"}
    status.DcAttributes = &common_proto.DataCenterAttributes{CreationDate:1234}


	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_HEARTBEAT,
		OpPayload: &common_proto.DCStream_DataCenter{DataCenter: &status},
	}

	log.Printf("messg %+v \n", &event)

	// publisher := micro2.NewPublisher("FromDcMgrToDcFacade")

	publisher := micro2.NewPublisher("FromDCFacadeToDCMgr") // appmgr to dcmgr

	publisher.Publish(&event)
}
