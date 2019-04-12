package main

import (
	"log"
	"time"

	"github.com/Ankr-network/dccn-common/pgrpc"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/config"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/dbservice"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/handler"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

// FIXME: delete if not used
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
	db, err := dbservice.New()
	if err != nil {
		//	log.Fatalln(err)
	}
	//defer db.Close()

	// init pgrpc
	if err := pgrpc.InitClient("50051" /*FIXME: hard code*/, nil); err != nil {
		log.Fatalln(err)
	}
	go handler.StartCollectStatus(db)
	time.Sleep(2 << 60)

	//// Initialise service
	//srv := grpc.NewService(
	//	micro.Name(ankr_default.DcMgrRegistryServerName),
	//)
	//srv.Init()
	//
	//// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	//srv.Server().Options().Broker.Connect()
	//if err := micro.RegisterSubscriber("dcMgrTaskDeploy", srv.Server(), &handler.Relay{}); err != nil {
	//	log.Fatalln(err)
	//}
	//
	//// Run srv
	//if err := srv.Run(); err != nil {
	//	log.Println(err.Error())
	//}
}
