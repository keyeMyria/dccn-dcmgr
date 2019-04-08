package handler

import (
	"golang.org/x/net/context"
	"encoding/json"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-dcmgr/app-dccn-dcmgr/db_service"
	"log"
)

type DcMgrAPIHandler struct {
	db dbservice.DBService
}

func NewAPIHandler(db dbservice.DBService) *DcMgrAPIHandler {
	handler := &DcMgrAPIHandler{
		db: db,
	}
	return handler
}

func (p *DcMgrAPIHandler) DataCenterList(
	ctx context.Context, req *common_proto.Empty)( *dcmgr.DataCenterListResponse, error) {
	//
	log.Println("api service receive DataCenterList from client")

	if list, err :=  p.db.GetAll(); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("DataCenterList successfully count: %d", len(*list))
		rsp := dcmgr.DataCenterListResponse{}
		rsp.DcList = *list
		return &rsp, nil
	}
}


func (p *DcMgrAPIHandler) DataCenterLeaderBoard(ctx context.Context, req *common_proto.Empty) ( *dcmgr.DataCenterLeaderBoardResponse, error){
	//rsp = & dcmgr.DataCenterLeaderBoardResponse{}
	dbList, err :=  p.db.GetAll();
	if err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	}


	list := make([]*dcmgr.DataCenterLeaderBoardDetail, 0)
	{
		detail := dcmgr.DataCenterLeaderBoardDetail{}
		detail.Name = "us_cloud"
		detail.Number = 99.81
		list = append(list, &detail)
	}

	{
		detail := dcmgr.DataCenterLeaderBoardDetail{}
		detail.Name = "asia_cloud"
		detail.Number = 97.71
		list = append(list, &detail)
	}

	{
		detail := dcmgr.DataCenterLeaderBoardDetail{}
		detail.Name = "europe_cloud"
		detail.Number = 96.89
		list = append(list, &detail)
	}


	for i := 0; i < len(*dbList); i++ {
		if i >= len(list) {
			break
		}
		dc := (*dbList)[i]
		list[i].Name = dc.Name
	}

	rsp := dcmgr.DataCenterLeaderBoardResponse{}

	rsp.List = list
	log.Printf("DataCenterLeaderBoard %+v", rsp.List)
	return &rsp, nil
}


func (p *DcMgrAPIHandler) NetworkInfo(ctx context.Context, req *common_proto.Empty)(*dcmgr.NetworkInfoResponse, error){
	rsp := dcmgr.NetworkInfoResponse{}
	rsp.UserCount = 299
	rsp.ContainerCount = 1342
	rsp.EnvironmentCount = 450
	rsp.HostCount = 137
	rsp.Traffic = p.calculateDCTraffic()
    return &rsp, nil
}

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

func (p *DcMgrAPIHandler) calculateDCTraffic() int32 {
	dbList, err :=  p.db.GetAll();
	if err == nil {
		totalCPU := 0
		usedCPU := 0
		for i := 0; i < len(*dbList); i++ {
			dc := (*dbList)[i]
			if dc.Status == common_proto.DCStatus_AVAILABLE {
				metrics := Metrics{}

				if err := json.Unmarshal([]byte(dc.DcHeartbeatReport.Metrics), &metrics); err != nil {
					log.Printf("metrics ")
				}else{
                    totalCPU += int(metrics.TotalCPU)
                    usedCPU  += int(metrics.UsedCPU)
				}
			}
		}

		if totalCPU == 0 {
			return 0 //  no dc available
		}else{
			rate := float64(usedCPU)/float64(totalCPU * 1000)
			if rate < 0.3 {   // only used 30%  it is light
				return 1
			}else if rate > 0.7 {  // used > 70%  it is heavy
				return 3
			}else{
				return 2       // median
			}




		}
	}


	return 0 //   no dc available
}
