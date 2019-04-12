package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/pgrpc"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/dbservice"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/handler"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config := micro2.LoadConfigFromEnv()
	config.Show()
}

func main() {
	db, err := dbservice.New()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	sendToDcMgr := micro2.NewPublisher("FromDCFacadeToDCMgr")
	if err := micro2.RegisterSubscriber("FromDcMgrToDcFacade", handler.New(sendToDcMgr)); err != nil {
		log.Fatalln(err)
	}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// init pgrpc
	if err := pgrpc.InitClient("50051" /*FIXME: hard code*/, nil); err != nil {
		log.Fatalln(err)
	}
	handler.StartCollectStatus(db)
}
