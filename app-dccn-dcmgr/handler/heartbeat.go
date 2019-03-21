package handler

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/db_service"
	"github.com/google/uuid"
	"log"
)

func (p *DcMgrHandler) UpdateDataCenter(ctx context.Context, dc *common_proto.DataCenter) error {
	// first update database
	//log.Printf("into updateDataCenter  : %v ", dc)
	center , err :=  p.db.GetByName(dc.Name)


    ip := dbservice.GetIP(ctx)
	//ip = "8.8.8.8"


	if center.Name == "" {
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


	return nil
}
