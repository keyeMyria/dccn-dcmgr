package scheduler

import (
	"container/heap"
	"encoding/json"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
	micro2 "github.com/Ankr-network/dccn-common/ankr-micro"
	"log"
)

type Metrics struct {
	TotalCPU     int64
	UsedCPU      int64
	TotalMemory  int64
	UsedMemory   int64
	TotalStorage int64
	UsedStorage  int64

	ImageCount    int64
	EndPointCount int64
	NetworkIO     int64 // No data
}

func DataCenterFilter(task *TaskRecord, db dbservice.DBService) []DataCenterRecord{
	list := make([]DataCenterRecord, 0)

	if len (task.Namespace.ClusterId) > 0 { // only one datacenter matching
		record := DataCenterRecord{}
		record.ID = task.Namespace.ClusterId
		record.Name = task.Namespace.ClusterName
		list = append(list, record)
		return list
	}


	dcs, _ := db.GetAvaliableList()
	micro2.Printf("avaliable datacenter %+v ", dcs)

	if len(*dcs) == 0 {
		return nil
	}



	for _, dc :=range *dcs {
		//check remain resource
		record := DataCenterRecord{}
		record.Name = dc.ClusterName
		record.ID = dc.DcId
		record.CPU = 2000
		record.Disk = 2000
		record.Memory = 3000
		//todo


		metrics := Metrics{}

		if err := json.Unmarshal([]byte(dc.DcHeartbeatReport.Metrics), &metrics); err != nil {
			log.Printf("datacenter metrics parse error ! ")
		} else {
			record.CPU =  metrics.TotalCPU * 1000 - metrics.UsedCPU
			record.Memory = metrics.TotalMemory - metrics.UsedMemory
			record.Disk = metrics.TotalStorage - metrics.UsedStorage
		}

		if record.CPU > int64(task.Namespace.NsCpuLimit) &&
			record.Memory > int64(task.Namespace.NsMemLimit) &&
			record.Disk > int64(task.Namespace.NsStorageLimit) {
			list = append(list, record)
			micro2.Printf("insert datacenter to DataCenterFilter list %s", dc.DcId)
		}else{
			//micro2.Printf(" datacenter %s not matching namespace %d %d %d   --> source requirement %d %d %d ", dc.Id,
			//	task.Namespace.CpuLimit, task.Namespace.MemLimit, task.Namespace.StorageLimit, record.CPU, record.Memory, record.Disk)
		}

        return list


	}




	return list
}


func DataCenterSelect(task *TaskRecord, db dbservice.DBService)string{
	dclist := DataCenterFilter(task, db)
	if dclist == nil {
		return ""
	}

	priorityQueue := make(DataCenterPriorityQueue, 0)

	for  _, dc := range dclist  {
		item := DataCenterRecordItem{}
		item.Record = dc
		item.Weight = 100  // TODO calcualate weight
		priorityQueue.Push(&item)
	}

    if len(priorityQueue) > 0 {
		item := heap.Pop(&priorityQueue).(*DataCenterRecordItem)
		return item.Record.ID
	}

	return ""
	// user pickup one best

}
