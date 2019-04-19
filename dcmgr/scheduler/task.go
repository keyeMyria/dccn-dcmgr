package scheduler

import "github.com/Ankr-network/dccn-common/protos/common"

type TaskRecord struct {
	ID           string
	Userid       string
	Name         string
	Image        string
	Namespace    *(common_proto.Namespace)
	Replica      int32
	Datacenterid string  // mongodb name is low field
	Status       common_proto.AppStatus // 1 new 2 running 3. done 4 cancelling 5.canceled 6. updating 7. updateFailed
	Msg  *common_proto.DCStream

}
