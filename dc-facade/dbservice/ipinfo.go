package dbservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type IPInfo struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Region   string `json:"region"`
}

func GetLatLng(ip string) (lat, lng, country string) {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s?token=05afd766593f88", ip))
	if err != nil {
		log.Printf("get ipinfo(%s) fail: %s", ip, err)
		return
	}
	defer resp.Body.Close()

	ipinfo := new(IPInfo)
	if err := json.NewDecoder(resp.Body).Decode(ipinfo); err != nil {
		log.Printf("parse ipinfo(%s) fail: %s", ip, err)
		return
	}

	values := strings.Split(ipinfo.Loc, ",")
	if len(values) != 2 {
		log.Printf("parse ipinfo(%s) fail, location: %s", ip, ipinfo.Loc)
		return
	}

	return values[0], values[1], ipinfo.Country
}
