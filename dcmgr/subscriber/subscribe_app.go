package subscriber

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
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
	if req.OpType == common_proto.DCOperation_NS_CREATE || req.OpType == common_proto.DCOperation_NS_UPDATE || req.OpType == common_proto.DCOperation_NS_CANCEL {
		p.dcFacadeDeploy.Publish(req)

	}else{
		task := req.GetAppDeployment()
		service := scheduler.GetSchedulerService()

		taskRecord := scheduler.TaskRecord{}
		taskRecord.Namespace = task.Namespace
		taskRecord.Msg = req
		service.AddTask(&taskRecord)

		//p.dcFacadeDeploy.Publish(req)
		micro2.Printf("add App to scheduler service %+v ", task)
	}




	return nil
}
