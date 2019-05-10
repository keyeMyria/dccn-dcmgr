package main

import (
	"log"

	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/pgrpc"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/handler"
	"github.com/Ankr-network/dccn-dcmgr/dc-facade/subscriber"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config := micro2.LoadConfigFromEnv()
	config.Show()
}

func main() {

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	sendToDcMgr := micro2.NewPublisher("FromDCFacadeToDCMgr")
	subscriber := subscriber.New(sendToDcMgr)
	heartbeat := handler.NewHeartBeat(sendToDcMgr)
	callback := handler.NewPgrpcHook(sendToDcMgr)

	if err := micro2.RegisterSubscriber("FromDcMgrToDcFacade", subscriber); err != nil {
		log.Fatalln(err)
	}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// init pgrpc
	if err := pgrpc.InitClient(micro2.GetConfig().Listen, nil, callback, handler.DialOpts()...); err != nil {
		log.Fatalln(err)
	}
	heartbeat.StartCollectStatus()
}
