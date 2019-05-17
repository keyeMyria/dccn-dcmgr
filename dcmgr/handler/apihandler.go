package handler

import (
	"encoding/json"
	certmanager "github.com/Ankr-network/dccn-common/cert"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"errors"
	"golang.org/x/net/context"
	"log"
	"time"
)


const (
	CA_CERT = `-----BEGIN CERTIFICATE-----
MIICKzCCAdKgAwIBAgIUW56lhwrEBMk7QKWYY/BZDAl3mLIwCgYIKoZIzj0EAwIw
dDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMQswCQYDVQQHEwJTRjEUMBIGA1UE
CRMLTUlTU0lPTiBTVC4xDjAMBgNVBBETBTk0MTA1MQ4wDAYDVQQKEwVIVUJDQTEV
MBMGA1UEAxMMbXlodWItY2EuY29tMB4XDTE5MDUxMjAxNDY1NVoXDTI5MDUxMjAx
NDY1NVowdDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMQswCQYDVQQHEwJTRjEU
MBIGA1UECRMLTUlTU0lPTiBTVC4xDjAMBgNVBBETBTk0MTA1MQ4wDAYDVQQKEwVI
VUJDQTEVMBMGA1UEAxMMbXlodWItY2EuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D
AQcDQgAEhIHiabGkTozh88/+TZrcAH0R5x5CrDN+Y4Czvt5AqqEfFXwU5Ihtt8a1
Pj87hqc6rQVYkwwc8Dgj3u60JgnZQaNCMEAwDgYDVR0PAQH/BAQDAgKEMB0GA1Ud
JQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MAoGCCqG
SM49BAMCA0cAMEQCIAYmfH43DBO8956XxOh+T3Dr/ijf0QmsIE1hb7jlh1MQAiAd
Lu3tEHLLmzQlal9FYxReurigU2kk0fngzbN7HAf8zQ==
-----END CERTIFICATE-----`

	CA_KEY = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIMIHHLuzL3k2/PAW6DhjGjATAGOx7cQK0lVK7AqW9C8GoAoGCCqGSM49
AwEHoUQDQgAEhIHiabGkTozh88/+TZrcAH0R5x5CrDN+Y4Czvt5AqqEfFXwU5Iht
t8a1Pj87hqc6rQVYkwwc8Dgj3u60JgnZQQ==
-----END EC PRIVATE KEY-----`
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
	ctx context.Context, req *common_proto.Empty) (*dcmgr.DataCenterListResponse, error) {
	//
	log.Println("api service receive DataCenterList from client")

	if list, err := p.db.GetAll(); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return nil, err
	} else {
		log.Printf("DataCenterList successfully count: %d", len(*list))
		rsp := dcmgr.DataCenterListResponse{}
		rsp.DcList = make([]*common_proto.DataCenterStatus, 0)

		for _, record :=range *list {
			cluster := p.GetClusterStatusFromClusterRecord(record)
			rsp.DcList = append(rsp.DcList, cluster)
		}
		return &rsp, nil
	}
}


//func (p *DcMgrAPIHandler) RegisterDataCenter(
//	ctx context.Context, req *dcmgr.RegisterDataCenterRequest) (*dcmgr.RegisterDataCenterResponse, error) {
	func (p *DcMgrAPIHandler) RegisterDataCenter(
		ctx context.Context, req *dcmgr.RegisterDataCenterRequest) (*dcmgr.RegisterDataCenterResponse, error) {

	uid := req.UserId
	dc, error := p.db.GetByUserID(uid)

	if error == nil {  // no found is ok
		return &dcmgr.RegisterDataCenterResponse{}, errors.New("datacenter exist")
	}

	log.Printf("dc %+v ------ %+v", dc, error)

	dataCenterId := uuid.New().String()


	rsp := dcmgr.RegisterDataCenterResponse{}
	cert, privateKey, err := GenerateClientCert(dataCenterId)

	if err == nil {
		rsp.CaCert = CA_CERT
		rsp.ClientCsrCert = cert
		rsp.ClientKey = privateKey

		if len(req.ClusterName) == 0 {
			req.ClusterName = "unknow"
		}

		dataCenter := dbservice.DataCenterRecord{}
		dataCenter.DcId = dataCenterId
		dataCenter.ClusterName = req.ClusterName
		dataCenter.DcStatus = common_proto.DCStatus_REGISTER
		dataCenter.GeoLocation = &common_proto.GeoLocation{}
		dataCenter.DcHeartbeatReport = &common_proto.DCHeartbeatReport{}
		dcAttributes := &common_proto.DataCenterAttributes{}
		now := time.Now().Unix()
		dcAttributes.CreationDate = &timestamp.Timestamp{Seconds: now}
		dcAttributes.LastModifiedDate = &timestamp.Timestamp{Seconds: now}
		dcAttributes.WalletAddress = ""
		dataCenter.DcAttributes = dcAttributes

		dataCenter.UserId = uid
		dataCenter.Clientcert = cert
		p.db.Create(&dataCenter)
	}

	return &rsp, nil
}


//func (p *DcMgrAPIHandler) RegisterDataCenter(
//	ctx context.Context, req *dcmgr.RegisterDataCenterRequest) (*dcmgr.RegisterDataCenterResponse, error) {
func (p *DcMgrAPIHandler) ResetDataCenter(
	ctx context.Context, req *dcmgr.RegisterDataCenterRequest) (*dcmgr.RegisterDataCenterResponse, error) {

	uid := req.UserId

	dc, error := p.db.GetByUserID(uid)

	if error != nil {  // no found
		return &dcmgr.RegisterDataCenterResponse{}, error
	}


	rsp := dcmgr.RegisterDataCenterResponse{}
	cert, privateKey, err := GenerateClientCert(dc.DcId)

	if err == nil {
		rsp.CaCert = CA_CERT
		rsp.ClientCsrCert = cert
		rsp.ClientKey = privateKey

		if len(req.ClusterName) == 0 {
			req.ClusterName = "unknow"
		}

		dataCenter := dbservice.DataCenterRecord{}
		dataCenter.ClusterName = req.ClusterName
		dataCenter.DcStatus = common_proto.DCStatus_REGISTER
		dataCenter.GeoLocation = &common_proto.GeoLocation{}
		dataCenter.DcHeartbeatReport = &common_proto.DCHeartbeatReport{}
		dataCenter.UserId = uid
		dataCenter.Clientcert = cert
		p.db.Reset(dc.DcId, &dataCenter)
	}

	return &rsp, nil
}

func (p *DcMgrAPIHandler)GetClusterStatusFromClusterRecord(record *dbservice.DataCenterRecord)*common_proto.DataCenterStatus{
	cluster := &common_proto.DataCenterStatus{}
	cluster.DcName = record.ClusterName
	cluster.DcId = record.DcId
	cluster.DcStatus = record.DcStatus
	cluster.DcHeartbeatReport = record.DcHeartbeatReport
	cluster.GeoLocation = record.GeoLocation
	cluster.DcAttributes = record.DcAttributes
	return cluster
}


func (p *DcMgrAPIHandler) MyDataCenter(
	ctx context.Context, req *dcmgr.MyDataCenterRequest) (*common_proto.DataCenterStatus, error) {
	uid := req.Uid
	record, error := p.db.GetByUserID(uid)
	cluster := p.GetClusterStatusFromClusterRecord(record)
	return cluster, error
}


func GenerateClientCert(datacenterId string)(string , string , error) {
	cert, privateKey, err := certmanager.GenerateEcdsaClientCert(datacenterId, CA_CERT, CA_KEY)

	log.Printf("GenerateClientCert-----> \n", err)
	if err == nil {
		log.Printf("\ncert: %s\n private key: %s\n", cert, privateKey)
		return cert, privateKey, err
	} else {
		return "", "", err
	}


}



func (p *DcMgrAPIHandler) NetworkInfo(ctx context.Context, req *common_proto.Empty) (*dcmgr.NetworkInfoResponse, error) {
	rsp := dcmgr.NetworkInfoResponse{}
	rsp.UserCount = 299
	rsp.ContainerCount = 1342
	rsp.NsCount = 450
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
	dbList, err := p.db.GetAll()
	if err == nil {
		totalCPU := 0
		usedCPU := 0
		for i := 0; i < len(*dbList); i++ {
			dc := (*dbList)[i]
			if dc.DcStatus == common_proto.DCStatus_AVAILABLE {
				metrics := Metrics{}

				if err := json.Unmarshal([]byte(dc.DcHeartbeatReport.Metrics), &metrics); err != nil {
					log.Printf("metrics ")
				} else {
					totalCPU += int(metrics.TotalCPU)
					usedCPU += int(metrics.UsedCPU)
				}
			}
		}

		if totalCPU == 0 {
			return 0 //  no dc available
		} else {
			rate := float64(usedCPU) / float64(totalCPU*1000)
			if rate < 0.3 { // only used 30%  it is light
				return 1
			} else if rate > 0.7 { // used > 70%  it is heavy
				return 3
			} else {
				return 2 // median
			}

		}
	}

	return 0 //   no dc available
}
