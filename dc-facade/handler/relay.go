package handler

import (
	"context"
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

func NewRelay(feedback *micro2.Publisher) *Relay {
	handler := &Relay{
		taskFeedback: feedback,
	}
	return handler
}



func (p *Relay)CreateApp(req *common_proto.DCStream) (err error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	app := req.GetAppDeployment()

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
		OpType:    common_proto.DCOperation_NS_CREATE,
		OpPayload: namespaceReport,
	}

	conn, err := pgrpc.Dial(app.Namespace.ClusterId)
	if err != nil {
		log.Println(err)
		appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
		return err  //todo use defer send msg for this case
	}
	defer conn.Close()

	resp, err := dcmgr.NewDCClient(conn).CreateApp(ctx, app)
	if err != nil {
		log.Println(err)
		appReport.AppReport.AppEvent = common_proto.AppEvent_LAUNCH_APP_FAILED
		return err //todo use defer send msg for this case
	} else {
		log.Printf("create app respone  %+v \n", resp)
		appReport.AppReport.AppEvent = resp.AppResult
		appReport.AppReport.Report = resp.Message
		namespaceReport.NsReport.NsEvent = resp.NsResult
	}


	p.taskFeedback.Publish(appEvent)
	p.taskFeedback.Publish(namespaceEvent)
	return nil

}


func (p *Relay)UpdateCancelDetailApp(req *common_proto.DCStream) (err error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	app := req.GetAppDeployment()

	appReport := &common_proto.DCStream_AppReport{
		AppReport: &common_proto.AppReport{
			AppDeployment: app,
		},
	}


	appEvent := &common_proto.DCStream{
		OpType:    req.OpType,
		OpPayload: appReport,
	}

	switch req.OpType {
	case common_proto.DCOperation_APP_UPDATE:
		conn, err := pgrpc.Dial(app.Namespace.ClusterId)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_UPDATE_APP_FAILED
			log.Println(err)
			return err //todo use defer send msg for this case
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).UpdateApp(ctx, app)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_UPDATE_APP_FAILED
			log.Println(err)
			return err
		} else {
			log.Printf("update app respone  %+v \n", resp)
			appReport.AppReport.AppEvent = resp.AppResult
			appReport.AppReport.Report = resp.Message
		}

	case common_proto.DCOperation_APP_CANCEL:
		conn, err := pgrpc.Dial(app.Namespace.ClusterId)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_CANCEL_APP_FAILED
			log.Println(err)
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).DeleteApp(ctx, app)
		if err != nil {
			appReport.AppReport.AppEvent = common_proto.AppEvent_CANCEL_APP_FAILED
			log.Println(err)
			return err
		} else {
			log.Printf("cancel app respone  %+v \n", resp)
			appReport.AppReport.AppEvent = resp.AppResult
			appReport.AppReport.Report = resp.Message
		}
	case common_proto.DCOperation_APP_DETAIL:
		conn, err := pgrpc.Dial(app.Namespace.ClusterId)
		if err != nil {
			log.Println(err)
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).Status(ctx, &common_proto.AppID{
			Id: app.Id,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		log.Printf("collect app detail of %s respone  %+v \n", app.Id, resp)
		appReport.AppReport.Detail = resp.Message

	default:
		log.Printf("process request error : request %+v \n", req)
		log.Println(ankr_default.ErrUnknown.Error())
		return ankr_default.ErrUnknown
	}

	p.taskFeedback.Publish(appEvent)
	return nil
}

func (p *Relay)NamespaceProcess(req *common_proto.DCStream) (err error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ns := req.GetNamespace()

	namespaceReport := &common_proto.DCStream_NsReport{
		NsReport: &common_proto.NamespaceReport{
			Namespace: ns,
		},
	}
	namespaceEvent := &common_proto.DCStream{
		OpType:    req.OpType,
		OpPayload: namespaceReport,
	}

	switch req.OpType {
	case common_proto.DCOperation_NS_CREATE:
		conn, err := pgrpc.Dial(ns.ClusterId)
		if err != nil {
			namespaceReport.NsReport.NsEvent = common_proto.NamespaceEvent_LAUNCH_NS_FAILED
			log.Println(err)
			return err  //todo use defer send msg for this case
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).CreateNamespace(ctx, ns)
		if err != nil {
			namespaceReport.NsReport.NsEvent = common_proto.NamespaceEvent_LAUNCH_NS_FAILED
			log.Println(err)
			return err
		}

		log.Printf("create namespace respone  %+v \n", resp)
		namespaceReport.NsReport.NsEvent = resp.NsResult

	case common_proto.DCOperation_NS_UPDATE:
		conn, err := pgrpc.Dial(ns.ClusterId)
		if err != nil {
			namespaceReport.NsReport.NsEvent = common_proto.NamespaceEvent_UPDATE_NS_FAILED
			log.Println(err)
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).UpdateNamespace(ctx, ns)
		if err != nil {
			namespaceReport.NsReport.NsEvent = common_proto.NamespaceEvent_LAUNCH_NS_FAILED
			log.Println(err)
			return err
		}

		log.Printf("update namespace respone  %+v \n", resp)
		namespaceReport.NsReport.NsEvent = resp.NsResult

	case common_proto.DCOperation_NS_CANCEL:
		conn, err := pgrpc.Dial(ns.ClusterId)
		if err != nil {
			namespaceReport.NsReport.NsEvent = common_proto.NamespaceEvent_CANCEL_NS_FAILED
			log.Println(err)
			return err
		}
		defer conn.Close()

		resp, err := dcmgr.NewDCClient(conn).DeleteNamespace(ctx, ns)
		if err != nil {
			namespaceReport.NsReport.NsEvent = common_proto.NamespaceEvent_LAUNCH_NS_FAILED
			log.Println(err)
			return err
		}

		log.Printf("delete namespace respone  %+v \n", resp)
		namespaceReport.NsReport.NsEvent = resp.NsResult

	default:
		log.Printf("process request error : request %+v \n", req)
		log.Println(ankr_default.ErrUnknown.Error())
		return ankr_default.ErrUnknown

	}

	p.taskFeedback.Publish(namespaceEvent)
	return nil
}
