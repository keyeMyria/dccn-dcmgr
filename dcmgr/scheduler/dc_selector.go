package scheduler

import "container/heap"

func DataCenterFilter(task *TaskRecord) []DataCenterRecord{
	list := make([]DataCenterRecord, 0)
	record := DataCenterRecord{}
	record.Name = "datacetner"
	record.ID = "12133"
	record.CPU = 123

	list = append(list, record)

	return list
}


func DataCenterSelect(task *TaskRecord)string{
	dclist := DataCenterFilter(task)

	priorityQueue := make(DataCenterPriorityQueue, 0)

	for  _, dc := range dclist  {
		item := DataCenterRecordItem{}
		item.Record = dc
		item.Weight = 100  // TODO calcualate weight
		priorityQueue.Push(&item)
	}

    if len(priorityQueue) > 0 {
		item := heap.Pop(&priorityQueue).(*DataCenterRecordItem)
		return item.Record.Name
	}

	return ""
	// user pickup one best

}
