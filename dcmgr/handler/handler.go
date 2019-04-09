package handler

import (
	"log"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/google/uuid"
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

	log.Printf("into updateTask from dc facade msg  : %v ", stream)
	p.taskFeedback.Publish(stream)
}

func (p *DcMgrHandler) UpdateDataCenter(dc *common_proto.DataCenter) error {
	// first update database
	//log.Printf("into updateDataCenter  : %v ", dc)
	center, err := p.db.GetByName(dc.Name)

	//ip := dbservice.GetIP(ctx)
	ip := "8.8.8.8"

	if center.Name == "" {
		// data center dose not exist, register it
		log.Printf("insert new datacenter  : %s  from ip : %s", dc.Name, ip)
		dc.Id = uuid.New().String()

		lat, lng, country := dbservice.GetLatLng(ip)
		dc.GeoLocation = &common_proto.GeoLocation{Lat: lat, Lng: lng, Country: country}

		if err = p.db.Create(dc); err != nil {
			log.Println(err.Error(), ", ", *dc)
			return err
		}
	} else {
		log.Printf("update datacenter by name : %s  ", center.Name)
		if err = p.db.Update(dc); err != nil {
			log.Println(err.Error())
			return err
		}
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
