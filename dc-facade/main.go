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

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	config := micro2.LoadConfigFromEnv()
	config.Show()

	if db, err = dbservice.New(); err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
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

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	sendToDcMgr := micro2.NewPublisher("FromDCFacadeToDCMgr")
	handler := handler.New(sendToDcMgr)

	if err := micro2.RegisterSubscriber("FromDcMgrToDcFacade", handler); err != nil {
		log.Fatalln(err)
	}

	forever := make(chan bool)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
