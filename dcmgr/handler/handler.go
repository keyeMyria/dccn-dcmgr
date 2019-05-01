package handler

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
)

type DcMgrHandler struct {
	db             dbservice.DBService
	taskFeedback   *micro2.Publisher       // sync task information with task manager
	DcStreamCaches *DataCenterStreamCaches // hold all data center as cache
}

func New(db dbservice.DBService, feedback *micro2.Publisher) *DcMgrHandler {
	handler := &DcMgrHandler{
		db:             db,
		taskFeedback:   feedback,
		DcStreamCaches: nil,
	}

	//	handler.DcStreamCaches.db = db
	return handler
}

//func (p *DcMgrHandler) ServerStream(
//	ctx context.Context, stream dcmgr.DCStreamer_ServerStreamStream) error {
//
//	log.Println("Debug into ServerStream")
//	for {
//		in, err := stream.Recv()
//		log.Println("Recv datacenter message")
//		if err == io.EOF {
//			log.Println("datacenter error eof ")
//			log.Println(err.Error())
//			return nil
//		}
//		if err != nil {
//			log.Println("datacenter error nil, dc may lost connection ")
//			log.Println(err.Error())
//			return err
//		}
//
//		switch in.OpType {
//		case common_proto.DCOperation_HEARTBEAT: // update data center in cache
//			if err := p.UpdateDataCenter(in.GetDataCenter()); err != nil {
//				log.Println(err.Error())
//			}
//		case common_proto.DCOperation_TASK_CREATE,
//			common_proto.DCOperation_TASK_UPDATE,
//			common_proto.DCOperation_TASK_CANCEL: // update task status
//			p.UpdateTask(in)
//		default:
//			log.Println(ankr_default.ErrUnknown.Error())
//		}
//	}
//}

func (p *DcMgrHandler) UpdateTask(stream *common_proto.DCStream) {

	micro2.Printf("update APP from dc facade msg  : %v ", stream)
	p.taskFeedback.Publish(stream)
}

func (p *DcMgrHandler) UpdateDataCenter(dc_status *common_proto.DataCenterStatus) error {
	// first update database
	//log.Printf("into updateDataCenter  : %v ", dc)
	if dc_status.Status == common_proto.DCStatus_UNAVAILABLE {
		p.db.UpdateStatus(dc_status.Id, common_proto.DCStatus_UNAVAILABLE)
		return nil
	}



	dc := new(common_proto.DataCenterStatus)
	dc.Name = dc_status.Name
	dc.Id = dc_status.Id
	dc.Status = dc_status.Status
	dc.DcHeartbeatReport = dc_status.DcHeartbeatReport
	dc.GeoLocation = dc_status.GeoLocation
	dc.DcAttributes = dc_status.DcAttributes


	datacenter, _ := p.db.Get(dc.Id)

	if len(datacenter.Id) == 0 {
		micro2.WriteLog("insert datacenter")
		p.db.Create(dc)
	}else{
		micro2.WriteLog("find datacenter, update datacenter \n")
		p.db.Update(dc)
	}

	return nil
}

func (p *DcMgrHandler) All() error {
	return nil
}

func (p *DcMgrHandler) Available() error {
	return nil
}

func (p *DcMgrHandler) Cleanup() {
	//if p.DcStreamCaches != nil {
	//	p.DcStreamCaches.Cleanup()
	//}
}