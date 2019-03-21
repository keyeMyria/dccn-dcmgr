package subscriber

import (
	"context"
	"github.com/micro/go-micro"
	"log"

	"github.com/Ankr-network/dccn-common/protos/common"

	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/handler"
)

type Subscriber struct {
	cache *handler.DataCenterStreamCaches
	dcFacadeDeploy micro.Publisher
}

func New(c *handler.DataCenterStreamCaches, dcFacadeDeploy micro.Publisher) *Subscriber {
	return &Subscriber{cache: c, dcFacadeDeploy: dcFacadeDeploy}
}

// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *Subscriber) HandlerDeploymentRequestFromTaskMgr(ctx context.Context, req *common_proto.DCStream) error {

	task := req.GetTask()
	//log.Printf("dc manager service(hub) HandlerDeployEvnetFromTaskMgr: Receive New Event: %+v", *task)
	//switch req.OpType {
	//case common_proto.DCOperation_TASK_CREATE,
	//	common_proto.DCOperation_TASK_CANCEL,
	//	common_proto.DCOperation_TASK_UPDATE:
	//	stream, err := p.cache.One(task.DataCenterName)
	//	if err != nil {
	//		log.Println(err.Error())
	//		return err
	//	}
	//	resp := &common_proto.DCStream{
	//		OpType:    req.OpType,
	//		OpPayload: &common_proto.DCStream_Task{Task: task}}
	//	if err := stream.Send(resp); err != nil {
	//		log.Println(err.Error())
	//		return err
	//	}
	//default:
	//	log.Println(ankr_default.ErrUnknown.Error())
	//	return ankr_default.ErrUnknown
	//}
	p.dcFacadeDeploy.Publish(context.TODO(), task)
	log.Printf("send message to DC Facade  %+v", *task)

	return nil
}
