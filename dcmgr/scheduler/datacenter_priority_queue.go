package scheduler

type DataCenterRecordItem struct {
	Record DataCenterRecord
	Weight int
	Index int // The index of the item in the heap.
}

type DataCenterPriorityQueue []*DataCenterRecordItem

func (pq DataCenterPriorityQueue) Len() int { return len(pq) }

func (pq DataCenterPriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest based on expiration number as the priority
	// The lower the expiry, the higher the priority
	return pq[i].Weight < pq[j].Weight
}

// We just implement the pre-defined function in interface of heap.

func (pq *DataCenterPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *DataCenterPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*DataCenterRecordItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq DataCenterPriorityQueue) Swap(i, j int) {
	if i < 0 || j < 0 || i >= pq.Len() || j >= pq.Len() {
		return
	}
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
