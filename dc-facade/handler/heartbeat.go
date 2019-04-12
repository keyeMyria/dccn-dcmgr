package handler

import (
	"context"
	"log"
	"time"

	"github.com/Ankr-network/dccn-common/pgrpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/dbservice"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func StartCollectStatus(db dbservice.DBService) {
	for range time.Tick(20 * time.Second) {
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

			// FIXME: transaction
			// update status into db
			if status.Id == "" {
				// data center dose not exist, register it
				lat, lng, country := dbservice.GetLatLng(key)
				status.GeoLocation = &common_proto.GeoLocation{Lat: lat, Lng: lng, Country: country}

				{ // init new dc id
					ts := time.Now().UTC().Unix()
					if _, err := dcmgr.NewDCClient(conn).InitDC(ctx, &common_proto.DataCenter{
						Id:   uuid.New().String(),
						Name: "mock_name",
						DcAttributes: &common_proto.DataCenterAttributes{
							CreationDate:     uint64(ts),
							LastModifiedDate: uint64(ts),
						},
					}); err != nil {
						log.Printf("init new datacenter fail: %s", err)
					}
				}

				log.Printf("add new datacenter: %s", status.Name)
				if err = db.Create(status); err != nil {
					log.Println(err.Error(), ", ", *status)
					return
				}

			} else {
				log.Printf("update datacenter by name : %s, ID: %s", status.Name, status.Id)
				if err = db.Update(status); err != nil {
					log.Println(err.Error())
					return
				}
			}
			return
		})
	}
}
