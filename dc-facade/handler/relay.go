package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/Ankr-network/dccn-common/pgrpc"
	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
)

type Relay struct{}

// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *Relay) HandlerDeploymentRequestFromTaskMgr(ctx context.Context, req *common_proto.DCStream) error {
	app := req.GetApp()
	if app == nil {
		return fmt.Errorf("invalid request data type: %T", req.OpPayload)
	}
	log.Printf("dc manager service(hub) HandlerDeployEvnetFromTaskMgr: Receive New Event: %+v", *app)

	switch req.OpType {
	case common_proto.DCOperation_TASK_CREATE:
		conn, err := pgrpc.Dial("app.DataCenterName")
		if err != nil {
			log.Println(err)
			return err
		}

		if _, err := dcmgr.NewDCClient(conn).CreateApp(ctx, app); err != nil {
			log.Println(err)
			return err
		}

	case common_proto.DCOperation_TASK_UPDATE:
		conn, err := pgrpc.Dial("app.DataCenterName")
		if err != nil {
			log.Println(err)
			return err
		}

		if _, err := dcmgr.NewDCClient(conn).UpdateApp(ctx, app); err != nil {
			log.Println(err)
			return err
		}

	case common_proto.DCOperation_TASK_CANCEL:
		conn, err := pgrpc.Dial("app.DataCenterName")
		if err != nil {
			log.Println(err)
			return err
		}

		if _, err := dcmgr.NewDCClient(conn).PurgeApp(ctx, &common_proto.AppID{
			Id: app.Id,
		}); err != nil {
			log.Println(err)
			return err
		}

	default:
		log.Println(ankr_default.ErrUnknown.Error())
		return ankr_default.ErrUnknown
	}

	log.Printf("send message to DataCenter  %+v", *app)
	return nil
}
