package subscriber

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/handler"
	"log"
)

type DCFacade struct {
	handler *handler.DcMgrHandler
}

func NewEventFromDCFacade(handler *handler.DcMgrHandler) *DCFacade {
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
	default:
	//case common_proto.DCOperation_APP_CREATE,
	//	common_proto.DCOperation_APP_UPDATE,
	//	common_proto.DCOperation_APP_CANCEL,
	//	common_proto.DCOperation_NS_CREATE
		p.handler.UpdateTask(in)

		//log.Println(ankr_default.ErrUnknown.Error())
	}

	micro2.Print("HandlerFeedBackFromDCFacade done ")

	return nil
}
