package handler

import (
	"encoding/json"
	certmanager "github.com/Ankr-network/dccn-common/cert"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"github.com/Ankr-network/dccn-dcmgr/dcmgr/db-service"
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
		rsp.DcList = *list
		return &rsp, nil
	}
}


//func (p *DcMgrAPIHandler) RegisterDataCenter(
//	ctx context.Context, req *dcmgr.RegisterDataCenterRequest) (*dcmgr.RegisterDataCenterResponse, error) {
	func (p *DcMgrAPIHandler) RegisterDataCenter(
		ctx context.Context, req *common_proto.Empty) (*dcmgr.NetworkInfoResponse, error) {


	//// todo check req.Uid
	//
	//dataCenterId := uuid.New().String()
	//
	//
	//rsp := dcmgr.RegisterDataCenterResponse{}
	//cert, privateKey, err := GenerateClientCert(dataCenterId)
	//
	//if err == nil {
	//	rsp.CaCert = CA_CERT
	//	rsp.ClientCsrCert = cert
	//	rsp.ClientKey = privateKey
	//
	//	dataCenter := &common_proto.DataCenterStatus{}
	//	dataCenter.Id = dataCenterId
	//	dataCenter.Name = "unknow"
	//	dataCenter.Status = common_proto.DCStatus_UNAVAILABLE
	//	dataCenter.GeoLocation = &common_proto.GeoLocation{}
	//	dataCenter.DcHeartbeatReport = &common_proto.DCHeartbeatReport{}
	//
	//	p.db.Create(dataCenter)
	//	p.db.UpdateClientCert(dataCenterId, cert)
	//}
	// todo after fix proxy problem
	rsp := dcmgr.NetworkInfoResponse{}

	return &rsp, nil
}
func GenerateClientCert(datacenterId string)(string , string , error) {
	cert, privateKey, err := certmanager.GenerateClientCert(datacenterId, CA_CERT, CA_KEY)

	if err == nil {
		log.Printf("\ncert: %s\n private key: %s\n", cert, privateKey)
		return cert, privateKey, err
	} else {
		return "", "", err
	}


}

func (p *DcMgrAPIHandler) DataCenterLeaderBoard(ctx context.Context, req *common_proto.Empty) (*dcmgr.DataCenterLeaderBoardResponse, error) {
	//rsp = & dcmgr.DataCenterLeaderBoardResponse{}
	dbList, err := p.db.GetAll()
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

func (p *DcMgrAPIHandler) NetworkInfo(ctx context.Context, req *common_proto.Empty) (*dcmgr.NetworkInfoResponse, error) {
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
	dbList, err := p.db.GetAll()
	if err == nil {
		totalCPU := 0
		usedCPU := 0
		for i := 0; i < len(*dbList); i++ {
			dc := (*dbList)[i]
			if dc.Status == common_proto.DCStatus_AVAILABLE {
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
