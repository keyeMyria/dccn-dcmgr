package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func main() {



	//app.Attributes = &common_proto.AppAttributes{}

	report := &common_proto.NamespaceReport{}
	namespace := common_proto.Namespace{
			Name:         "test_ns1",
			CpuLimit:     300,
			MemLimit:     500,
			StorageLimit: 10,
		}

	report.Namespace = &namespace

	report.NsEvent = common_proto.NamespaceEvent_LAUNCH_NS
	report.NsStatus = common_proto.NamespaceStatus_NS_LAUNCHING


	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_NS_CREATE,
		OpPayload: &common_proto.DCStream_NsReport{NsReport: report},
	}

	log.Printf("messg %+v \n", &event)

	// publisher := micro2.NewPublisher("FromDcMgrToDcFacade")

	publisher := micro2.NewPublisher("topic.deploy.app") // appmgr to dcmgr

	publisher.Publish(&event)
}
