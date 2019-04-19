package scheduler

import (
	"container/heap"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
	"log"
	"time"
)

// each queue for one datacenter. datacenter pickuping task bases one his priority rules

var service *SchedulerService

var started = false

var LoopInterval = 1   // for debug it is 1 second, for production  10 second

type SchedulerService struct {
	queues    map[string]*PriorityQueue // this is task queue for publish to dc_facade
	publisher *micro2.Publisher
	db dbservice.DBService
}

func GetSchedulerService() *SchedulerService {
	if service == nil {
		log.Printf("SchedulerService is nil, not start properly")
	}
	return service
}

func New(p *micro2.Publisher, db dbservice.DBService) *SchedulerService {
	service = new(SchedulerService)
	service.queues = make(map[string]*PriorityQueue)
	service.publisher = p
	service.db = db
	return service
}

func (s *SchedulerService) AddTask(task *TaskRecord) {
	dcID := DataCenterSelect(task, s.db)
	if len(dcID) > 0 {
		queue := s.GetTaskPriorityQueue(dcID)
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
		for k, v := range s.queues {
			if (len(*v)) > 0 {
				item := heap.Pop(v).(*TaskRecordItem)
				s.SendTaskToDataCenter(k, item.Task)
			}

		}

		time.Sleep(time.Duration(LoopInterval) * time.Second)
	}

}

func (s *SchedulerService) SendTaskToDataCenter(datacenterID string, task *TaskRecord) {
	// deploy to dc_facade
	taskCreateMsg := task.Msg
	appDeployment := taskCreateMsg.GetAppDeployment()
	if appDeployment.Namespace.ClusterId == datacenterID {
		s.publisher.Publish(task.Msg)  // no need add clusterid
	}else{
		appDeployment.Namespace.ClusterId = datacenterID
		appDeployment.Namespace.Name = s.getDatacenterName(datacenterID)
		event := common_proto.DCStream{
			OpType:    common_proto.DCOperation_APP_CREATE,
			OpPayload: &common_proto.DCStream_AppDeployment{AppDeployment: appDeployment},
		}
		s.publisher.Publish(&event)

	}

	log.Printf("SendTaskToDataCenter  task id %s , cluser id %s  \n", appDeployment.Id , appDeployment.Namespace.ClusterId)
}

func (s *SchedulerService)getDatacenterName(datacenterId string) string {
	 record, err := s.db.Get(datacenterId)
	 if err == nil {
         return ""
	 }else{
	 	return record.Name
	 }

}

func (s *SchedulerService) Start() {
	if started == false {
		go s.LoopForSchedule()
	}
}
