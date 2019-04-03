package subscriber

import (
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/handler"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/micro"
	"log"
)

type Subscriber struct {
	cache *handler.DataCenterStreamCaches
	dcFacadeDeploy *micro2.Publisher
}

func New(c *handler.DataCenterStreamCaches, dcFacadeDeploy *micro2.Publisher) *Subscriber {
	return &Subscriber{cache: c, dcFacadeDeploy: dcFacadeDeploy}
}

func (p *Subscriber) HandlerDeploymentRequestFromTaskMgr(req *common_proto.DCStream) error {
   //   debug.PrintStack()
	task := req.GetTask()
	p.dcFacadeDeploy.Publish(req)
	log.Printf("send message to DC Facade  %+v", task)

	return nil
}


