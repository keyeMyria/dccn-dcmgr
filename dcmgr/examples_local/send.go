package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func main() {
	app := common_proto.App{}
	app.Name = "workpress_test"
	app.ChartDetail = &common_proto.ChartDetail{
		Repo:    "stable",
		Name:    "workpress",
		Version: "5.7.1",
	}
	app.NamespaceData = &common_proto.App_Namespace{
		Namespace: &common_proto.Namespace{
			Name:         "test_ns1",
			CpuLimit:     300,
			MemLimit:     500,
			StorageLimit: 10,
		},
	}

	app.Attributes = &common_proto.AppAttributes{}

	appDeployment := &common_proto.AppDeployment{}
	appDeployment.Id = "1111"
	appDeployment.Name = "2222"

	appDeployment.ChartDetail = app.ChartDetail
	appDeployment.Uid = "xxxxx"

	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_APP_CREATE,
		OpPayload: &common_proto.DCStream_AppDeployment{AppDeployment: appDeployment},
	}

	log.Printf("messg %+v \n", &event)

	// publisher := micro2.NewPublisher("FromDcMgrToDcFacade")

	publisher := micro2.NewPublisher("topic.deploy.app") // appmgr to dcmgr

	publisher.Publish(&event)
}
