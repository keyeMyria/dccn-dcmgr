package subscriber

import (
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/handler"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"log"
)

type Subscriber struct {
	relay *(handler.Relay)
}

func New(feedback *micro2.Publisher) *Subscriber {
	subscriber := Subscriber{}
	handler := handler.NewRelay(feedback)
	subscriber.relay = handler

	return &subscriber
}


// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *Subscriber) HandlerDeploymentRequestFromDcMgr(req *common_proto.DCStream) (err error) {
	log.Printf("dc manager service(hub) HandlerDeployEvnetFromDcMgr: Receive New Event: %+v %+v", req)

	switch req.OpType {
	case common_proto.DCOperation_APP_CREATE:
		return p.relay.CreateApp(req)
	case common_proto.DCOperation_APP_UPDATE, common_proto.DCOperation_APP_CANCEL, common_proto.DCOperation_APP_DETAIL:
		return p.relay.UpdateCancelDetailApp(req)
	case common_proto.DCOperation_NS_CREATE, common_proto.DCOperation_NS_UPDATE, common_proto.DCOperation_NS_CANCEL:
		return p.relay.NamespaceProcess(req)
	case common_proto.DCOperation_HEARTBEAT:
		fallthrough
	default:
		log.Printf("process request error : request %+v \n", req)
		log.Println(ankr_default.ErrUnknown.Error())
		return ankr_default.ErrUnknown
	}

	log.Printf("send message to DcMgr")
	return nil
}