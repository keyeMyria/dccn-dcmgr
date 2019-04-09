package scheduler

import "github.com/Ankr-network/dccn-common/protos/common"

type TaskRecord struct {
	ID           string
	Userid       string
	Name         string
	Image        string
	Datacenter   string
	//Type         common_proto.TaskType
	Replica      int32
	Datacenterid string  // mongodb name is low field
	Status       common_proto.AppStatus // 1 new 2 running 3. done 4 cancelling 5.canceled 6. updating 7. updateFailed
	Hidden       bool
	Schedule     string
	Last_modified_date uint64
	Creation_date uint64
	Msg  *common_proto.DCStream

}
