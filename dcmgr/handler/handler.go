package handler

import (
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
)

type DcMgrHandler struct {
	db             dbservice.DBService
	taskFeedback   *micro2.Publisher       // sync task information with task manager
}

func New(db dbservice.DBService, feedback *micro2.Publisher) *DcMgrHandler {
	handler := &DcMgrHandler{
		db:             db,
		taskFeedback:   feedback,
	}

	//	handler.DcStreamCaches.db = db
	return handler
}


func (p *DcMgrHandler) UpdateTask(stream *common_proto.DCStream) {

	micro2.Printf("update APP from dc facade msg  : %v ", stream)
	p.taskFeedback.Publish(stream)
}

func (p *DcMgrHandler) UpdateDataCenter(dc_status *common_proto.DataCenterStatus) error {
	// first update database
	//log.Printf("into updateDataCenter  : %v ", dc)
	if dc_status.DcStatus == common_proto.DCStatus_UNAVAILABLE {
		p.db.UpdateStatus(dc_status.DcId, common_proto.DCStatus_UNAVAILABLE)
		return nil
	}


	dc := new(common_proto.DataCenterStatus)
	//dc.Name = dc_status.Name
	//dc.Id = dc_status.Id
	dc.DcStatus = dc_status.DcStatus
	dc.DcHeartbeatReport = dc_status.DcHeartbeatReport
	dc.GeoLocation = dc_status.GeoLocation
	dc.DcAttributes = dc_status.DcAttributes


	datacenter, _ := p.db.Get(dc.DcId)

	if len(datacenter.DcId) == 0 {
		micro2.WriteLog("error update datacenter failed for datacenterid does not exist")
		//p.db.Create(dc)
	}else{
		micro2.WriteLog("find datacenter, update datacenter")
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
