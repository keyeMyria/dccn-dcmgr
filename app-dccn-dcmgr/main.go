package main

import (
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/handler"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/subscriber"
	"log"

	"github.com/micro/go-micro"

	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/config"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/db_service"
	"github.com/micro/go-grpc"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	conf config.Config
	db   dbservice.DBService
	err  error
)

func main() {
	Init()

	if db, err = dbservice.New(conf.DB); err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	startHandler()
}

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if conf, err = config.Load(); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Load config %+v\n", conf)
}

func startHandler() {
	srv := grpc.NewService(
		micro.Name(ankr_default.DcMgrRegistryServerName),
	)

	// Initialise service
	srv.Init()

	// Register Task Handler
	log.Println("Registering DC Handler")
	dcClient := handler.NewAPIHandler(db)
	if err := dcmgr.RegisterDCAPIHandler(srv.Server(), dcClient); err != nil {
		log.Fatal(err.Error())
	}
	//
	//// Dc Manager register handler
	// New Publisher to deploy new task action.
	//New Publisher to deploy new task action.
	taskFeedback := micro.NewPublisher(ankr_default.MQFeedbackTask, srv.Client())
	dcFacadeDeploy := micro.NewPublisher("dcMgrTaskDeploy", srv.Client())

	dcHandler := handler.New(db, taskFeedback)

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	opt := srv.Server().Options()
	opt.Broker.Connect()
	if err := micro.RegisterSubscriber(ankr_default.MQDeployTask, srv.Server(), subscriber.New(dcHandler.DcStreamCaches, dcFacadeDeploy)); err != nil {
		log.Fatal(err.Error())
	}




	//from
	if err := micro.RegisterSubscriber("FromDCFacadeToDCMgr", srv.Server(), subscriber.NewEventFromDCFacade(dcHandler.DcStreamCaches, dcHandler)); err != nil {
		log.Fatal(err.Error())
	}


	//defer dcHandler.Cleanup()

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
