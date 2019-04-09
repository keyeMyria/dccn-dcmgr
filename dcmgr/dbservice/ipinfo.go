package dbservice

import (
	"encoding/json"
	"context"
	"fmt"
	"google.golang.org/grpc/peer"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)


type IPInfo struct {
	IP       string
	Hostname string
	City     string
	Country  string
	Loc      string
	Org      string
}

func GetLatLng(ip string) (string, string, string)  {
	reader := strings.NewReader(``)

	url := "https://ipinfo.io/"+ip+"?token=05afd766593f88"
	fmt.Print(url)
	request, err := http.NewRequest("GET", url, reader)
	if err != nil {
		fmt.Print(err.Error())
		return "", "", ""
	}
	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		fmt.Print("ipinfo error :" + err.Error())
		return "", "", "US"
	}

	contents, err := ioutil.ReadAll(resp.Body)

	var ipinfo IPInfo

	err = json.Unmarshal(contents, &ipinfo)

	values :=  strings.Split(ipinfo.Loc, ",")

	if err != nil  || len(values) < 2{
		fmt.Printf("GetLatLng eror >>>> for ip %s \n", ip)
		return "", "", "US"
	}


	log :=fmt.Sprintf("loc lat: %s lng: %s \n", values[0], values[1])
	fmt.Print(log)
	return values[0], values[1], ipinfo.Country
}

func GetIP(ctx context.Context) string {
	log.Printf("context %+v \n", ctx)
	pr, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Print("failed to get peer from ctx " )
	}
	if pr.Addr == net.Addr(nil) {
		fmt.Print("failed to get peer address")
	}

	values :=  strings.Split(pr.Addr.String(), ":")
	fmt.Print(" >>  IP : " + values[0])
	return values[0]
}
