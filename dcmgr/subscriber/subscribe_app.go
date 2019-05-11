package subscriber

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/scheduler"
)

type Subscriber struct {
	dcFacadeDeploy *micro2.Publisher
}

func New(dcFacadeDeploy *micro2.Publisher) *Subscriber {
	return &Subscriber{dcFacadeDeploy: dcFacadeDeploy}
}

func (p *Subscriber) HandlerDeploymentRequestFromTaskMgr(req *common_proto.DCStream) error {
	//   debug.PrintStack()
	if req.OpType == common_proto.DCOperation_NS_CREATE || req.OpType == common_proto.DCOperation_NS_UPDATE || req.OpType == common_proto.DCOperation_NS_CANCEL {

		if req.OpType == common_proto.DCOperation_NS_CREATE {
			ns := req.GetNamespace()
			if len(ns.ClusterId) == 0 { //need set clusterid
				service := scheduler.GetSchedulerService()
				service.AddNamespace(req)
                return nil
			}
		}

		ns := req.GetNamespace()
		if len (ns.ClusterName) == 0 {
			service := scheduler.GetSchedulerService()
			service.AddNamespace(req)
			return nil
		}
		p.dcFacadeDeploy.Publish(req)

		micro2.Printf("send namespace create/update/cancel to  datacenter %+v ", req)

	}else{

		if req.OpType == common_proto.DCOperation_APP_CREATE {
			task := req.GetAppDeployment()
			service := scheduler.GetSchedulerService()

			taskRecord := scheduler.TaskRecord{}
			taskRecord.Namespace = task.Namespace
			taskRecord.Msg = req
			service.AddTask(&taskRecord)

			//p.dcFacadeDeploy.Publish(req)
			micro2.Printf("add App to scheduler service %+v ", task)
		}else if req.OpType == common_proto.DCOperation_APP_UPDATE || req.OpType == common_proto.DCOperation_APP_CANCEL || req.OpType == common_proto.DCOperation_APP_DETAIL {
			p.dcFacadeDeploy.Publish(req)
			micro2.Printf("send app update/cancel/detail to  datacenter %+v ", req)
		}else{
			micro2.Printf("error in HandlerDeploymentRequestFromTaskMgr %+v", req)
		}




	}




	return nil
}
