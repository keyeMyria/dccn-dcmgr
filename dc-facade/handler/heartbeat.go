package handler

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/dbservice"
	"google.golang.org/grpc"
)

var db dbservice.DBService

func CollectStatus() {
	for range time.Tick(20 * time.Second) {
		pgrpc.Each(func(key string, conn *grpc.ClientConn, err error) (stop bool) {
			// handle dial error
			stop = false
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

			// FIXME: transaction
			// update status into db
			center, err := db.GetByName(status.Name)
			if err != nil {
				log.Println(err)
				return
			}
			if center.Name == "" {
				// data center dose not exist, register it
				host, _, err := net.SplitHostPort(key)
				if err != nil {
					log.Println(err)
					return
				}
				lat, lng, country := dbservice.GetLatLng(host)
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
