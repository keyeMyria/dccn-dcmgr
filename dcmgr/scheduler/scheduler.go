package scheduler

import (
	"container/heap"
	"context"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	app "github.com/Ankr-network/dccn-common/protos/appmgr/v1/grpc"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
	"google.golang.org/grpc"
	"log"
	"time"
)

// each queue for one datacenter. datacenter pickuping task bases one his priority rules

var service *SchedulerService

var started = false

var LoopInterval = 1 // for debug it is 1 second, for production  10 second

type SchedulerService struct {
	queues    map[string]*PriorityQueue // this is task queue for publish to dc_facade
	publisher *micro2.Publisher
	db        dbservice.DBService
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
		item.Weight = UserRunningApplicationOnDataCenter(task.Userid, dcID) //todo  wight =    user used resouce  / dc total resouce   quesion:
		queue.Push(&item)
	} else {
		//todo send failed create msg to appmgr

		log.Printf("can not find data center, add task failed\n")
	}

}

func (s *SchedulerService) AddNamespace(req *common_proto.DCStream) {
	dcs, _ := s.db.GetAvaliableList()
	if len(*dcs) > 0 {
		dc := (*dcs)[0]
		ns := req.GetNamespace()
		ns.ClusterId = dc.DcId
		ns.ClusterName = dc.DcName
		event := common_proto.DCStream{
			OpType:    common_proto.DCOperation_NS_CREATE,
			OpPayload: &common_proto.DCStream_Namespace{Namespace: ns},
		}
		log.Printf("send namespace create msg to %s , %+v", ns.NsName, ns)
		s.publisher.Publish(&event)

	} else {
		log.Printf("not database avalible for namespace ")
	}

}

func UserRunningApplicationOnDataCenter(user_id string, cluster_id string) int {
	count := 65
	var addr = "appmgr:50051"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("connect to %s error \n", addr)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)

	appClient := app.NewAppMgrClient(conn)

	if rsp, err := appClient.AppCount(context.Background(), &app.AppCountRequest{UserId: user_id, ClusterId: cluster_id}); err != nil {
		log.Printf("error when call AppCount %s \n" + err.Error())
	} else {
		count = int(rsp.AppCount)
	}

	log.Printf("UserRunningApplicationOnDataCenter  %d \n", count)
	// call api get
	return 1000 - count
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
		s.publisher.Publish(task.Msg) // no need add clusterid
	} else {
		appDeployment.Namespace.ClusterId = datacenterID
		appDeployment.Namespace.ClusterName = s.getDatacenterName(datacenterID)
		event := common_proto.DCStream{
			OpType:    common_proto.DCOperation_APP_CREATE,
			OpPayload: &common_proto.DCStream_AppDeployment{AppDeployment: appDeployment},
		}
		s.publisher.Publish(&event)

	}

	log.Printf("SendTaskToDataCenter  task id %s , cluster id %s cluster name: %s \n", appDeployment.AppId,
		appDeployment.Namespace.ClusterId, appDeployment.Namespace.ClusterName)
}

func (s *SchedulerService) getDatacenterName(datacenterId string) string {
	record, err := s.db.Get(datacenterId)
	log.Printf("hello %+v \n", record)
	if err != nil {
		return ""
	} else {
		return record.DcName
	}

}

func (s *SchedulerService) Start() {
	if started == false {
		go s.LoopForSchedule()
	}
}
