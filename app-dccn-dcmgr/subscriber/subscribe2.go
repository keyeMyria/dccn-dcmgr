package subscriber

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos"
	"log"

	"github.com/Ankr-network/dccn-common/protos/common"

	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/handler"
)

type DCFacade struct {
	handler *handler.DcMgrHandler
}

func NewEventFromDCFacade(c *handler.DataCenterStreamCaches, handler *handler.DcMgrHandler) *DCFacade {
	return &DCFacade{handler : handler}
}

// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *DCFacade) HandlerFeedBackFromDCFacade(ctx context.Context, req *common_proto.DCStream) error {

	in := req
	switch in.OpType {
	case common_proto.DCOperation_HEARTBEAT: // update data center in cache

		if err := p.handler.UpdateDataCenter(ctx, in.GetDataCenter()); err != nil {
			log.Println(err.Error())
		}
	case common_proto.DCOperation_TASK_CREATE,
		common_proto.DCOperation_TASK_UPDATE,
		common_proto.DCOperation_TASK_CANCEL: // update task status
		if err := p.handler.UpdateTask(in); err != nil {
			log.Println(err.Error())
		}
	default:
		log.Println(ankr_default.ErrUnknown.Error())
	}


	log.Printf("HandlerFeedBackFromDCFacade done ")

	return nil
}
