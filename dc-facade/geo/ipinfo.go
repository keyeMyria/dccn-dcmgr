package geo

import (
	"encoding/json"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
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

func ReadIPinfo(ip string)(lat, lng, country string) {
	{
		db, err := geoip2.Open("GeoIP2-City.mmdb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		// If you are using strings that may be invalid, check that ip is not nil
		ip := net.ParseIP(ip)
		record, err := db.City(ip)
		if err != nil {
			log.Printf("error when ReadIPinfo")
			return "", "", ""
		}

		latitude := fmt.Sprintf("%f", record.Location.Latitude)
		longtitude := fmt.Sprintf("%f", record.Location.Longitude)

		return latitude, longtitude, record.Country.IsoCode

		/*
		fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
		fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
		fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
		fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
		fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
		fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
		*/
	}
}
