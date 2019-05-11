package handler

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

const (
	HUB_CERT = `-----BEGIN CERTIFICATE-----
MIICJzCCAc6gAwIBAgIUQNK8zuB47TrjMK/9apa4+ODmGP8wCgYIKoZIzj0EAwIw
dDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMQswCQYDVQQHEwJTRjEUMBIGA1UE
CRMLTUlTU0lPTiBTVC4xDjAMBgNVBBETBTk0MTA1MQ4wDAYDVQQKEwVIVUJDQTEV
MBMGA1UEAxMMbXlodWItY2EuY29tMB4XDTE5MDUxMjAxNDY1NVoXDTI5MDUxMjAx
NDY1NVowfTELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMQswCQYDVQQHEwJTRjEU
MBIGA1UECRMLTUlTU0lPTiBTVC4xDjAMBgNVBBETBTk0MTA1MRMwEQYDVQQKEwpE
YXRhQ2VudGVyMRkwFwYDVQQDExBteWRhdGFjZW50ZXIuY29tMFkwEwYHKoZIzj0C
AQYIKoZIzj0DAQcDQgAEM49mdr428vS5+uHc0wjJBqyQ5n8d0QLra97C40uaEw94
l6RWjMOGbQfHGg6YbZzQ6Zc0qIxf7xu+RX//sTmqCaM1MDMwDgYDVR0PAQH/BAQD
AgeAMBMGA1UdJQQMMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwCgYIKoZIzj0E
AwIDRwAwRAIgUxRoNWAjjyvTmnzU8c8s02g0wZURKGo76kh9LNVXcp4CIBAvaZ5u
Y88YwWeiSVJNBDC6MIcgPLAM4YuLvNjP6M6W
-----END CERTIFICATE-----`

	HUB_KEY = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIAHFNZ8+2UnV72fsnUciUAoHYiBKY+FO7IZoT2TPMUUaoAoGCCqGSM49
AwEHoUQDQgAEM49mdr428vS5+uHc0wjJBqyQ5n8d0QLra97C40uaEw94l6RWjMOG
bQfHGg6YbZzQ6Zc0qIxf7xu+RX//sTmqCQ==
-----END EC PRIVATE KEY-----`

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
)

func DialOpts() []grpc.DialOption {
	cert, err := tls.X509KeyPair([]byte(HUB_CERT), []byte(HUB_KEY))
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}

	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM([]byte(CA_CERT))
	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		ClientCAs:          certPool,
		InsecureSkipVerify: true, // FIXME: turn to false if cert is dynamic sign
	}
	transportCreds := credentials.NewTLS(&tlsConfig)

	return []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5,
			Timeout:             20,
			PermitWithoutStream: true,
		}),
	}
}
