package main

import (
	"context"
	"log"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	testCommon "github.com/Ankr-network/dccn-hub/dcmgr/examples/common"

	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

// send events using the publisher
func sendEv(taskId string, p micro.Publisher) {

	// create new event
	ev := common_proto.Event{
		EventType: common_proto.Operation_TASK_CANCEL,
		OpMessage: &common_proto.Event_TaskFeedback{TaskFeedback: &common_proto.TaskFeedback{
			TaskId: taskId,
			Status: common_proto.TaskStatus_CANCEL_FAILED,
		}},
	}

	log.Printf("publishing %+v\n", ev)

	// publish an event
	if err := p.Publish(context.Background(), &ev); err != nil {
		log.Fatalf("error publishing %v\n", err)
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// create a service
	service := grpc.NewService()

	// parse command line
	service.Init()

	// create publisher
	pub := micro.NewPublisher(ankr_default.MQDeployTask, service.Client())

	task := testCommon.MockTasks()[0]
	if err := pub.Publish(context.Background(), &common_proto.Event{EventType: common_proto.Operation_TASK_CREATE, OpMessage: &common_proto.Event_Task{Task: &task}}); err != nil {
		log.Fatal(err.Error())
	}
}
