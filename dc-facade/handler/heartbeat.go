package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/dbservice"
	"google.golang.org/grpc"
)

func StartCollectStatus(db dbservice.DBService) {
	for range time.Tick(3 * time.Second) {
		pgrpc.Each(func(key string, conn *grpc.ClientConn, err error) (stop bool) {
			// handle dial error
			stop = true
			if err != nil {
				log.Println(err)
				return
			}

			// collect status(heartbeat)
			ctx := context.Background()
			status, err := dcmgr.NewDCClient(conn).Status(ctx, &common_proto.Empty{})
			if err != nil {
				log.Println(err)
				return
			}

			// FIXME: enable mongodb
			data, _ := json.MarshalIndent(status, "", "    ")
			log.Printf("%s", data)
			return

			// FIXME: transaction
			// update status into db
			center, err := db.GetByName(status.Name)
			if err != nil {
				log.Println(err)
				return
			}
			if center.Name == "" {
				// data center dose not exist, register it
				lat, lng, country := dbservice.GetLatLng(key)
				status.GeoLocation = &common_proto.GeoLocation{Lat: lat, Lng: lng, Country: country}

				log.Printf("add new datacenter: %s", status.Name)
				if err = db.Create(status); err != nil {
					log.Println(err.Error(), ", ", *status)
					return
				}
			} else {
				log.Printf("update datacenter by name : %s  ", center.Name)
				if err = db.Update(status); err != nil {
					log.Println(err.Error())
					return
				}
			}
			return
		})
	}
}
