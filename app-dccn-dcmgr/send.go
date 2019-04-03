package main

import (
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/micro"
)

func main() {

	//
	//e := micro2.Event{}
	//e.Name = "ssssss"
	//
	//publisher := micro2.NewPublisher(ankr_default.MQDeployTask)

	task := common_proto.Task{}
	//task.DataCenterName = "datacenter_tokyo"
	task.Name = "task"
	task.Type = common_proto.TaskType_DEPLOYMENT
	task.Attributes = &common_proto.TaskAttributes{}
	task.Attributes.Replica = 1
	//task.ChartName = "deploymentzys"
	//task.ChartVer = "0.0.1"
	//task.Uid = "f6045dc7-d46a-40c1-b637-d4c75b6213fa"
	t := common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image:"nginx:1.12"}}
	task.TypeData = &t

	task.Id = "570d72bd-32bc-41dc-9540-bb05cf1fae28"


	event := common_proto.DCStream{
		OpType: common_proto.DCOperation_TASK_CREATE,
		OpPayload: &common_proto.DCStream_Task{Task: &task},
	}
	//
	//fmt.Printf("the event sring %+v \n", event)
	////publisher.Publish(event)
	//b, _ := json.Marshal(event)
	//fmt.Printf("the event sring ----> %+v \n", string(b))
	//msg := &common_proto.DCStream{}
	//json.Unmarshal(b, msg)
	//
	//fmt.Printf("\n\n Unmarshal  %+v \n", *msg)

	micro2.Send("test", &event)




}
