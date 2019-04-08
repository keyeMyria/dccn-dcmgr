package subscriber

import (
	"log"

	"github.com/Ankr-network/dccn-common/protos/common"
	micro2 "github.com/Ankr-network/dccn-dcmgr/dcmgr/ankr-micro"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/handler"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/scheduler"
)

type Subscriber struct {
	cache          *handler.DataCenterStreamCaches
	dcFacadeDeploy *micro2.Publisher
}

func New(c *handler.DataCenterStreamCaches, dcFacadeDeploy *micro2.Publisher) *Subscriber {
	return &Subscriber{cache: c, dcFacadeDeploy: dcFacadeDeploy}
}

func (p *Subscriber) HandlerDeploymentRequestFromTaskMgr(req *common_proto.DCStream) error {
	//   debug.PrintStack()
	task := req.GetTask()
	service := scheduler.GetSchedulerService()
	taskRecord := scheduler.TaskRecord{}
	taskRecord.Msg = req
	service.AddTask(&taskRecord)

	//p.dcFacadeDeploy.Publish(req)
	log.Printf("add AddTask to scheduler service %+v  \n", task)

	return nil
}
