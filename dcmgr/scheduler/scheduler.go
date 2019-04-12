package scheduler

import (
	"container/heap"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"log"
	"time"
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
	// deploy to dc_facade
	s.publisher.Publish(task.Msg)
	log.Printf("SendTaskToDataCenter %+v\n", task.Msg)
}

func (s *SchedulerService) Start() {
	if started == false {
		go s.LoopForSchedule()
	}
}
