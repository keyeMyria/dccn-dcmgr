package handler

import (
	"encoding/json"
	certmanager "github.com/Ankr-network/dccn-common/cert"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
	"github.com/google/uuid"
	"errors"
	"golang.org/x/net/context"
	"log"

)


const (
	CA_CERT = `
-----BEGIN CERTIFICATE-----
MIIDvTCCAqWgAwIBAgICBnUwDQYJKoZIhvcNAQELBQAwYjERMA8GA1UEBhMIU2hh
bWJhbGExDzANBgNVBAgTBnNoYW0gMDELMAkGA1UEBxMCVUExDDAKBgNVBAoTA1pF
TjENMAsGA1UECxMET20gMDESMBAGA1UEAxMJbG9jYWxob3N0MB4XDTE5MDUwMTIx
MzQ1M1oXDTI5MDUwMTIxMzQ1M1owYjERMA8GA1UEBhMIU2hhbWJhbGExDzANBgNV
BAgTBnNoYW0gMDELMAkGA1UEBxMCVUExDDAKBgNVBAoTA1pFTjENMAsGA1UECxME
T20gMDESMBAGA1UEAxMJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAuohTX8Iclq0bpmZK5rpJCTIAT/Gxz+r8Wk/thJOrxcHaIm/LoGch
ZCOCVow5gobnUOzfBSzv6ny5UcD6bUBTIPWMhaK1sxmvHnnIbQWKBmaRzRYtOe28
2dB11TpM4qmpLkEigpBeAMf9Yb40dg9xi/MGkRI6Ky+lD7FJaqAoJHkpCSd2FO1K
StInoDnm1obqUGQE7AwgzG35y/j7gTdiUPhqr+CoH1/7fqvIMp7dDu52KXs4u8Gb
IyCa/dOb3BfOyNYZ/wLDOux6CxS2J74LON56n5Cc+2/15ym5bfUuSzCoyqQxyhvr
lXwupObFHr2XuYGJZMHk6R1drQ6ha40ziQIDAQABo30wezAOBgNVHQ8BAf8EBAMC
AoQwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMA8GA1UdEwEB/wQFMAMB
Af8wHQYDVR0OBBYEFPd9Sz1ooSLqo3i0MyU+9zd96HFXMBoGA1UdEQQTMBGCCWxv
Y2FsaG9zdIcEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEAcNqbR1Pj2xTDzkRr3M56
H70pHy2gGgdMOldzVwfOWErapmyI/KLjXCfNCq6iZzMsa0jnlg8YVQVxSgdfahC8
/aA+U4X2n3yGEN2phIuSmKUaWpt7c0OyulN05L/yo9sTg6XIcvV1Uy2nv9KqgsUp
pMu/rh6TUsVwlWLDhvhFeMsO5I07h4sTwuxNSBB02/VVkp0t7nZvhg5ZbnZZy82w
W6aNCauQVQlU29C7VzwBzL7wKotBuPCpC944NG5MKmls551dK/GgusX/eMDQuXO4
UNGpHZXtIYn+BvR8ILH1dBHEnFADl3l1j7lc/nlQ3OAUpGC0U+2sL7YG9XkRAmR+
Rg==
-----END CERTIFICATE-----`

	CA_KEY = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAuohTX8Iclq0bpmZK5rpJCTIAT/Gxz+r8Wk/thJOrxcHaIm/L
oGchZCOCVow5gobnUOzfBSzv6ny5UcD6bUBTIPWMhaK1sxmvHnnIbQWKBmaRzRYt
Oe282dB11TpM4qmpLkEigpBeAMf9Yb40dg9xi/MGkRI6Ky+lD7FJaqAoJHkpCSd2
FO1KStInoDnm1obqUGQE7AwgzG35y/j7gTdiUPhqr+CoH1/7fqvIMp7dDu52KXs4
u8GbIyCa/dOb3BfOyNYZ/wLDOux6CxS2J74LON56n5Cc+2/15ym5bfUuSzCoyqQx
yhvrlXwupObFHr2XuYGJZMHk6R1drQ6ha40ziQIDAQABAoIBAEkqM9m5n9ESMWhB
c1uw8fjjXe/9k4tEVptuAnIgIh49fTxIsYxHJ3fJ3dPoyJ4EIDND1e6Hw8ssBNym
XxP/SRwCdI3uVmrbxi6kAhOROqRsEwBolHDGaW7eL3nllkbJ2YxFDC4+RkD0MNTn
8FfmktkcCBVbGunZlFrlZTCnhVdYazg+jmOlld1CHkwqO1vZArBvKy6K92EzpRk7
3IRaDIqKt3ALCDEM39t+B7mdzVeIm2y1k3Iq62s2CKKk7ui/icNhbDtHPWg1ZDrJ
jl+VxFBhyarnu7inzyttnd8QkGi3ii4vQf3d2KTpPOJldcA8pnIIkGWW8+tCPH5M
zjxiIAECgYEA5XVLkA3pjARK3wgVgpyUSANeIdHD4wmxpUlbUQaweIYhm3Yud3EC
FMtSWIOCUdZ8g1760gtkIgmzYr62XNco42Fv+scZZQhBNfSSiWoqeKq6j9Mp69jx
tNAMwZ12Y5DzaD/3FQU01Cfn0ntQRI66aaOKDWlXozGIa5Fulbk/ygECgYEA0Bvq
+EUT8OrZJidxiW53dWmL2zLCZW9KeoENgV3hRSsHOqxByw+iJkYQIKy54fzkZ6dq
tZrPi3td2EwLhiXfVsUoqQps204v5HZLRT9lw18iSJAySF1aqdMB3gOy69uFn5xe
W1GkPkJ6qx6CpJisLGfq1RdyDKKB7U1tDS2wGYkCgYEAhSIWUqHP1R6Udm2RVXQW
EOZrUoIL/wob2YQDiLKx52wjybi7Yy/dfkUuJQ9AqM0i93I/Y2makqlAPNXcp2dr
YOqi90VX9afhdjXOZA6GT/b3QgXKN/5q13czP49mJoTuxZj/emHH8iSpPBWyT+Tk
QfDSY8+wOo690XPTTunqGgECgYEArOCkPv3TZO0S2rklfg9AOU8mmT7/chgTfNS8
DV2Zh0YJSVpThYZFIxpMx3f1KqBUdS8EXDxwcORYvxfc8uF/OKur7VD1wPCgpF8I
hEv4E2ZyKmlu++JhMHZTNMVJ2tiPllnloGKf2ACNup0r1ePmEzV4RPCnE4vj9ue8
0ZfElFECgYEAyDj/gw7l9V6cxUuafxgfeEaCcDYuC2Rc7Z8WY+TyM6sASUR3lyvi
+uEkI1VIIvagla/J18gjALEei7hVS/ywTG2YSrzdDwo5H5Kf+D+W3EcP35og4VZa
4V36uydYuiG1Z4X5YWHsr6S3e0hg79ZgGO7Rug07JE5R3neS8x5jqiA=
-----END RSA PRIVATE KEY-----`
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
	uid := req.UserId
	record, error := p.db.GetByUserID(uid)
	cluster := p.GetClusterStatusFromClusterRecord(record)
	return cluster, error
}


func GenerateClientCert(datacenterId string)(string , string , error) {
	cert, privateKey, err := certmanager.GenerateClientCert(datacenterId, CA_CERT, CA_KEY)

	if err == nil {
		//log.Printf("\ncert: %s\n private key: %s\n", cert, privateKey)
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
