package main

import (
	"fmt"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/micro"
	"golang.org/x/net/context"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/db_service"
)

type HandleNew struct {
	id int32
}

func (h HandleNew)Handle(req *common_proto.DCStream){
   fmt.Printf("handle %+v \n", req)
}


type GRPCHandler struct {
	db int
}
func NewGRPCAPIHandler()  * GRPCHandler{
	handler := &GRPCHandler{
	}
	return handler
}

func (p *GRPCHandler) DataCenterList(
	ctx context.Context, req *common_proto.Empty)( *dcmgr.DataCenterListResponse, error) {

	return &dcmgr.DataCenterListResponse{} , nil
}



func (p *GRPCHandler) DataCenterLeaderBoard(ctx context.Context, req *common_proto.Empty)( *dcmgr.DataCenterLeaderBoardResponse, error) {
	//rsp = & dcmgr.DataCenterLeaderBoardResponse{}
	return nil,nil
}


func (p *GRPCHandler) NetworkInfo(ctx context.Context, req *common_proto.Empty,
	)( *dcmgr.NetworkInfoResponse, error){

	return nil, nil
}


func main() {


	db, _ := dbservice.New()
	db.GetAll()


	handler := HandleNew{}
	micro2.RegisterSubscriber("test", handler)

	dcClient := NewGRPCAPIHandler()
	service := micro2.NewService()
	dcmgr.RegisterDCAPIServer(service.GetServer(), dcClient)
	service.Start()



}
