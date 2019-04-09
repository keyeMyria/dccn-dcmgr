package main

import (
	"log"

	"github.com/Ankr-network/dccn-common/pgrpc"
	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/config"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/dbservice"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/handler"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var conf config.Config

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var err error
	if conf, err = config.Load(); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Load config %+v\n", conf)
}

func main() {
	// connect to mongo
	db, err := dbservice.New(conf.DB)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// init pgrpc
	if pgrpc.InitClient(":50051", nil); err != nil {
		log.Fatalln(err)
	}

	// Initialise service
	srv := grpc.NewService(
		micro.Name(ankr_default.DcMgrRegistryServerName),
	)
	srv.Init()

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	srv.Server().Options().Broker.Connect()
	if err := micro.RegisterSubscriber("dcMgrTaskDeploy", srv.Server(), &handler.Relay{}); err != nil {
		log.Fatalln(err)
	}

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
//
//// Init starts handler to listen.
//func Init() {
//	log.SetFlags(log.Lshortfile | log.LstdFlags)
//
//	if conf, err = config.Load(); err != nil {
//		log.Fatal(err.Error())
//	}
//	log.Printf("Load config %+v\n", conf)
//}
//
//func startHandler() {
//	srv := grpc.NewService(
//		micro.Name(ankr_default.DcMgrRegistryServerName),
//	)
//
//	// Initialise service
//	srv.Init()
//
//	// Register Task Handler
//
//	//Dc Manager register handler
//	//New Publisher to deploy new task action.
//	taskFeedback := micro.NewPublisher("FromDCFacadeToDCMgr", srv.Client())
//
//	dcHandler := handler.New(db, taskFeedback)
//
//	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
//	opt := srv.Server().Options()
//	opt.Broker.Connect()
//	if err := micro.RegisterSubscriber("dcMgrTaskDeploy", srv.Server(), subscriber.New(dcHandler.DcStreamCaches)); err != nil {
//		log.Fatal(err.Error())
//	}
//
//	// Register Dc Manager Handler
//	if err := dcmgr.RegisterDCStreamerHandler(srv.Server(), dcHandler); err != nil {
//		log.Fatal(err.Error())
//	}
//
//
//	//defer dcHandler.Cleanup()
//
//	// Run srv
//	if err := srv.Run(); err != nil {
//		log.Println(err.Error())
//	}
//}
