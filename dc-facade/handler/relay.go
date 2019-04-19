package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/pgrpc"
	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
)

type Relay struct {
	taskFeedback *micro2.Publisher
}

func New(feedback *micro2.Publisher) *Relay {
	handler := &Relay{
		taskFeedback: feedback,
	}
	return handler
}

func (p *Relay) sendTestMsg(msg *common_proto.DCStream) {
	report := common_proto.AppReport{}
	app := msg.GetAppDeployment()
	app.Namespace.ClusterName = "dc"
	app.Namespace.ClusterId = "2e8556cb-17dd-4584-9adc-a58d36f92ce5"
	app.Namespace.CreationDate = uint64(time.Now().Second())
	report.AppDeployment = app
	report.Report = "this is a fake msg for test"
	report.AppStatus = common_proto.AppStatus_APP_RUNNING

	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_APP_CREATE,
		OpPayload: &common_proto.DCStream_AppReport{AppReport: &report},
	}

	p.taskFeedback.Publish(&event)
	log.Printf("SendTaskToDCMgr  %+v\n", event)
}

// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *Relay) HandlerDeploymentRequestFromDcMgr(req *common_proto.DCStream) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	app := req.GetAppDeployment()
	if app == nil {
		return fmt.Errorf("invalid request data type: %T", req.OpPayload)
	}
	log.Printf("dc manager service(hub) HandlerDeployEvnetFromDcMgr: Receive New Event: %+v", *app)

	//p.sendTestMsg(req)  this is test message
	appReport := &common_proto.DCStream_AppReport{
		AppReport: &common_proto.AppReport{
			AppDeployment: app,
		},
	}
	appEvent := &common_proto.DCStream{
		OpType:    req.OpType,
		OpPayload: appReport,
	}


	namespaceReport := &common_proto.DCStream_NsReport{
		NsReport: &common_proto.NamespaceReport{
			Namespace: app.Namespace,
		},
	}

	namespaceEvent := &common_proto.DCStream{
		OpType:    req.OpType,
		OpPayload: namespaceReport,
	}
	defer func() {
		p.taskFeedback.Publish(namespaceEvent)
		p.taskFeedback.Publish(appEvent)
	}()

	switch req.OpType {
	case common_proto.DCOperation_APP_CREATE:
		conn, err := pgrpc.Dial(app.Namespace.ClusterId)
		if err != nil {
			log.Println(err)
			appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).CreateApp(ctx, app)
		if err != nil {
			log.Println(err)
			appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
			return err
		}else{
			appReport.AppReport.AppEvent =  resp.AppResult
			appReport.AppReport.Report =  resp.Message
			namespaceReport.NsReport.NsEvent = resp.NsResult
		}

        /*
		if resp.NsResult != common_proto.NamespaceEvent_LAUNCH_NS_SUCCEED {
			appEvent.OpPayload = &common_proto.DCStream_NsReport{
				NsReport: &common_proto.NamespaceReport{
					Namespace: app.Namespace,
					NsEvent:   resp.NsResult,
				},
			}
		}
		appReport.AppReport = toReport(resp)

        */

	case common_proto.DCOperation_APP_UPDATE:
		conn, err := pgrpc.Dial(app.Namespace.ClusterId)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
			log.Println(err)
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).UpdateApp(ctx, app)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
			log.Println(err)
			return err
		}else{
		appReport.AppReport.AppEvent =  resp.AppResult
		appReport.AppReport.Report =  resp.Message
		namespaceReport.NsReport.NsEvent = resp.NsResult
	}

		//if resp.NsResult != common_proto.NamespaceEvent_LAUNCH_NS_SUCCEED {
		//	appEvent.OpPayload = &common_proto.DCStream_NsReport{
		//		NsReport: &common_proto.NamespaceReport{
		//			Namespace: app.Namespace,
		//			NsEvent:   resp.NsResult,
		//		},
		//	}
		//}
		//appReport.AppReport = toReport(resp)

	case common_proto.DCOperation_APP_CANCEL:
		conn, err := pgrpc.Dial(app.Namespace.ClusterId)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
			log.Println(err)
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).DeleteApp(ctx, app)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
			log.Println(err)
			return err
		}else{
			appReport.AppReport.AppEvent =  resp.AppResult
			appReport.AppReport.Report =  resp.Message
			namespaceReport.NsReport.NsEvent = resp.NsResult
		}

		//if resp.NsResult != common_proto.NamespaceEvent_LAUNCH_NS_SUCCEED {
		//	appEvent.OpPayload = &common_proto.DCStream_NsReport{
		//		NsReport: &common_proto.NamespaceReport{
		//			Namespace: app.Namespace,
		//			NsEvent:   resp.NsResult,
		//		},
		//	}
		//}
		//appReport.AppReport = toReport(resp)

	default:
		log.Println(ankr_default.ErrUnknown.Error())
		return ankr_default.ErrUnknown
	}

	// set event with new data
	appEvent = &common_proto.DCStream{
		OpType:    req.OpType,
		OpPayload: appReport,
	}

	namespaceEvent = &common_proto.DCStream{
		OpType:    req.OpType,
		OpPayload: namespaceReport,
	}


	log.Printf("send message to DataCenter  %+v", *app)
	return nil
}

func toReport(resp *common_proto.AppResponce) *common_proto.AppReport {
	return &common_proto.AppReport{
		Report:   resp.Error,
		AppEvent: resp.AppResult,
		Detail:   resp.Message,
	}
}
