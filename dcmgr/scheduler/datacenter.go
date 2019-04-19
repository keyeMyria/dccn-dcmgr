package scheduler

type DataCenterRecord struct {
	ID           string
	Name         string
	CPU          int64
	Memory       int64
	Disk         int64
    Price        float64
}
