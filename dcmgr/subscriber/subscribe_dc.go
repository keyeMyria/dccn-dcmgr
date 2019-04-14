package subscriber

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/handler"
	"log"
)

type DCFacade struct {
	handler *handler.DcMgrHandler
}

func NewEventFromDCFacade(c *handler.DataCenterStreamCaches, handler *handler.DcMgrHandler) *DCFacade {
	return &DCFacade{handler: handler}
}

func (p *DCFacade) HandlerFeedBackFromDCFacade(req *common_proto.DCStream) error {

	micro2.Printf("recevie messge from dc-facade %+v \n ", req)

	in := req
	switch in.OpType {
	case common_proto.DCOperation_HEARTBEAT: // update data center in cache

		if err := p.handler.UpdateDataCenter(in.GetDataCenter()); err != nil {
			log.Println(err.Error())
		}
	case common_proto.DCOperation_APP_CREATE,
		common_proto.DCOperation_APP_UPDATE,
		common_proto.DCOperation_APP_CANCEL: // update task status
		p.handler.UpdateTask(in)
	default:
		log.Println(ankr_default.ErrUnknown.Error())
	}

	micro2.Print("HandlerFeedBackFromDCFacade done ")

	return nil
}
