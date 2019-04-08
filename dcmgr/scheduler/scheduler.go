package scheduler

import (
	"container/heap"
	"fmt"
	"log"
	"time"

	"github.com/Ankr-network/dccn-common/protos/common"
	micro2 "github.com/Ankr-network/dccn-dcmgr/dcmgr/ankr-micro"
)

// each queue for one datacenter. datacenter pickuping task bases one his priority rules

var service *SchedulerService

var started = false

var LoopInterval = 10

type SchedulerService struct {
	queues    map[string]*PriorityQueue // this is task queue for publish to dc_facade
	publisher *micro2.Publisher
}

func GetSchedulerService() *SchedulerService {
	if service == nil {
		log.Printf("SchedulerService is nil, not start properly")
	}
	return service
}

func New(p *micro2.Publisher) *SchedulerService {
	service = new(SchedulerService)
	service.queues = make(map[string]*PriorityQueue)
	service.publisher = p
	return service
}

func (s *SchedulerService) AddTask(task *TaskRecord) {
	dc := DataCenterSelect(task)
	if len(dc) > 0 {
		queue := s.GetTaskPriorityQueue(dc)
		item := TaskRecordItem{}
		item.Task = task
		item.Weight = 100
		queue.Push(&item)
	} else {
		log.Printf("can not find data center, add task failed\n")
	}

}

func (s *SchedulerService) GetTaskPriorityQueue(datacenter string) *PriorityQueue {
	queue, _ := s.queues[datacenter]
	if queue == nil {
		taskQueue := make(PriorityQueue, 0)
		s.queues[datacenter] = &taskQueue
	}

	return s.queues[datacenter]

}

func (s *SchedulerService) LoopForSchedule() {
	for {
		//	log.Printf("LoopForSchedule >>> \n")
		for k, v := range s.queues {
			if (len(*v)) > 0 {
				item := heap.Pop(v).(*TaskRecordItem)
				s.SendTaskToDataCenter(k, item.Task)
			}

		}

		time.Sleep(time.Duration(LoopInterval) * time.Second)
	}

}

func (s *SchedulerService) SendTaskToDataCenter(datacenter string, task *TaskRecord) {
	// TODO update  task  fields (status and datacenter) in database
	// deploy to dc_facade
	log.Printf("SendTaskToDataCenter %+v\n", task.Msg)
	s.publisher.Publish(task.Msg)
	//send2(s.publisher, task.Msg)

}

func (s *SchedulerService) Start() {
	if started == false {
		go s.LoopForSchedule()
	}
}

func send2(publisher *micro2.Publisher, task2 *common_proto.Task) {

	e := micro2.Event{}
	e.Name = "ssssss"

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
	t := common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image: "nginx:1.12"}}
	task.TypeData = &t

	task.Id = "570d72bd-32bc-41dc-9540-bb05cf1fae28"

	event := common_proto.DCStream{
		OpType:    common_proto.DCOperation_TASK_CREATE,
		OpPayload: &common_proto.DCStream_Task{Task: &task},
	}
	////
	//fmt.Printf("the event sring %+v \n", event)
	////publisher.Publish(event)
	//b, _ := json.Marshal(event)
	//fmt.Printf("the event sring ----> %+v \n", string(b))
	//msg := &common_proto.DCStream{}
	//json.Unmarshal(b, msg)
	//
	fmt.Printf("task2 %+v \n", task2)
	fmt.Printf("event %+v \n", event)

	//micro2.Send("dcMgrTaskDeploy", &event)
	publisher.Publish(&event)

}
