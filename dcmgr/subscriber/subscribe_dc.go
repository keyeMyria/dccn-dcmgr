package subscriber

import (
	"github.com/Ankr-network/dccn-common/protos"
	"log"

	"github.com/Ankr-network/dccn-common/protos/common"

	"github.com/Ankr-network/dccn-dcmgr/dcmgr/handler"
)

type DCFacade struct {
	handler *handler.DcMgrHandler
}

func NewEventFromDCFacade(c *handler.DataCenterStreamCaches, handler *handler.DcMgrHandler) *DCFacade {
	return &DCFacade{handler : handler}
}

func (p *DCFacade) HandlerFeedBackFromDCFacade(req *common_proto.DCStream) error {

	in := req
	switch in.OpType {
	case common_proto.DCOperation_HEARTBEAT: // update data center in cache

		if err := p.handler.UpdateDataCenter(in.GetDataCenter()); err != nil {
			log.Println(err.Error())
		}
	case common_proto.DCOperation_TASK_CREATE,
		common_proto.DCOperation_TASK_UPDATE,
		common_proto.DCOperation_TASK_CANCEL: // update task status
		p.handler.UpdateTask(in)
	default:
		log.Println(ankr_default.ErrUnknown.Error())
	}


	log.Printf("HandlerFeedBackFromDCFacade done ")

	return nil
}
