package handler

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/common"
	"log"
)

func (p *DcMgrHandler) UpdateTask(stream *common_proto.DCStream) error {

	log.Printf("into updateTask from dc facade msg  : %v ", stream)
	return p.taskFeedback.Publish(context.TODO(), stream)
}
