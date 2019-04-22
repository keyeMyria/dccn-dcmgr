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

func (p *Relay) StartCollectStatus() {
	for range time.Tick(20 * time.Second) {
		pgrpc.Each(func(key string, conn *grpc.ClientConn, err error) {
			// handle dial error
			if err != nil {
				log.Println(err)
				return
			}
			defer conn.Close()

			// collect status(heartbeat)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			log.Println("collecting status of key ", key)
			status, err := dcmgr.NewDCClient(conn).Overview(ctx, &common_proto.Empty{})
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
					status.Id = uuid.New().String()
					status.Name = "mock_name"
					ts := uint64(time.Now().UTC().Unix())

					log.Println("initing dc of key ", key)
					if _, err := dcmgr.NewDCClient(conn).InitDC(ctx, &common_proto.DataCenter{
						Id:   status.Id,
						Name: status.Name,
						DcAttributes: &common_proto.DataCenterAttributes{
							CreationDate:     ts,
							LastModifiedDate: ts,
						},
					}); err != nil {
						log.Printf("init new datacenter fail: %s", err)
						return
					}

					log.Printf("added new datacenter: %s", status.Name)
				}
			}

			if key != status.Id {
				log.Printf("alias %s into %s", key, status.Id)
				pgrpc.Alias(key, status.Id, true)
			}

			event := common_proto.DCStream{
				OpType: common_proto.DCOperation_HEARTBEAT,
				OpPayload: &common_proto.DCStream_DataCenter{
					DataCenter: status,
				},
			}

			log.Println("publishing status of key ", key)
			p.taskFeedback.Publish(&event)
		})

		log.Printf("heartbeat finish")
	}
}
