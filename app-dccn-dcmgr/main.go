package main

import (
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	//"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/handler"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/subscriber"
	//"github.com/micro/go-micro"
	"log"

	//"github.com/micro/go-micro"

	"github.com/Ankr-network/dccn-common/protos"
	//"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/config"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/db_service"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/micro"
	//	"github.com/micro/go-grpc"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	db   dbservice.DBService
	err  error
)

func main() {
	Init()

	if db, err = dbservice.New(); err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	startHandler()
}

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	config := micro2.LoadConfigFromEnv()

	config.Show()
}

func startHandler() {

	//New Publisher to deploy new task action.
	taskFeedback := micro2.NewPublisher(ankr_default.MQFeedbackTask )
	dcFacadeDeploy := micro2.NewPublisher("dcMgrTaskDeploy")

	dcHandler := handler.New(db, taskFeedback)


	if err := micro2.RegisterSubscriber(ankr_default.MQDeployTask,  subscriber.New(dcHandler.DcStreamCaches, dcFacadeDeploy)); err != nil {
		log.Fatal(err.Error())
	}

	//from
	if err := micro2.RegisterSubscriber("FromDCFacadeToDCMgr", subscriber.NewEventFromDCFacade(dcHandler.DcStreamCaches, dcHandler)); err != nil {
		log.Fatal(err.Error())
	}

	service := micro2.NewService()


	dcClient := handler.NewAPIHandler(db)
	dcmgr.RegisterDCAPIServer(service.GetServer(), dcClient)
	service.Start()


}
