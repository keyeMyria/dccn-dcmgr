package handler

import (
	"log"
	"context"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
	"github.com/google/uuid"
	dbservice "github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/db_service"
)

func (p *DcMgrHandler) updateDataCenter(ctx context.Context, dc *common_proto.DataCenter, stream dcmgr.DCStreamer_ServerStreamStream) error {
	// first update database
	//log.Printf("into updateDataCenter  : %v ", dc)
	center , err :=  p.db.GetByName(dc.Name)


	if center.Name == "" {
		log.Printf("updateDataCenter start get IP")
		ip := dbservice.GetIP(ctx)
		log.Printf("updateDataCenter end get IP")
		//ip = "8.8.8.8"

		// data center dose not exist, register it
		log.Printf("insert new datacenter  : %s  from ip : %s", dc.Name, ip)
		dc.Id = uuid.New().String()

		lat, lng, country := dbservice.GetLatLng(ip)
		dc.GeoLocation = &common_proto.GeoLocation{Lat:lat, Lng:lng, Country:country}

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

	// then update stream
	log.Printf("update new data center stream: %s ", dc.Name)
	p.DcStreamCaches.Add(dc, stream)

	return nil
}
