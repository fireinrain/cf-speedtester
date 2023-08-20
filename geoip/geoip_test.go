package geoip

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"testing"
)

func TestGeoIPGitRepo_GetRepoLatestTag(t *testing.T) {
	repo := GlobalGeoIPClient

	fmt.Println(repo)

}

func TestGeoIPGitRepo_GEOIPInfo(t *testing.T) {
	repo := GlobalGeoIPClient
	db, err := geoip2.Open(repo.GeoIPDbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP("146.70.175.116")
	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	//if len(record.Subdivisions) > 0 {
	//	fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	//}
	//fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	//fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	//fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)

}

func TestName(t *testing.T) {
	var age []string
	fmt.Println(age == nil)
}
