package scheduler

import (
	"database/sql"
	"dewietl/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type ValidatorGeoData struct {
	Host          string  `json:"host"`
	IP            string  `json:"ip"`
	RDNS          string  `json:"rdns"`
	ASN           int     `json:"asn"`
	ISP           string  `json:"isp"`
	CountryName   string  `json:"country_name"`
	CountryCode   string  `json:"country_code"`
	RegionName    string  `json:"region_name"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	PostalCode    string  `json:"postal_code"`
	ContinentName string  `json:"continent_name"`
	ContinentCode string  `json:"continent_code"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Metro_code    int     `json:"metro_code"`
	Timezone      string  `json:"timezone"`
	Datetime      string  `json:"datetime"`
}

type Validator struct {
	Address string `json:"address"`
	IP      string `json:"ip"`
}

func updateValidatorGeoData() {

	for {

		var address, listen_addrs sql.NullString

		rows, err := database.DB.Query("SELECT address, listen_addrs FROM validator_status WHERE listen_addrs IS NOT NULL")
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&address, &listen_addrs)
			if err != nil {
				log.Printf("[ERROR] getAllHotspots(): %v", err)
			}

			if address.Valid == true && listen_addrs.Valid == true {

				validatorIP := getValidatorIP(listen_addrs.String)

				geoInformation := getGeoData(validatorIP)

				geoString, err := json.Marshal(geoInformation)
				if err != nil {
					fmt.Println(err)
				}

				initialString := "INSERT INTO validator_isp (address, isp, geo_data, host, ip, rdns, asn, country_name, country_code, region_name, region_code, city, postal_code, continent_name, continent_code, latitude, longitude, metro_code, timezone, datetime) VALUES "
				queryString := fmt.Sprintf(" ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v') ", address.String, geoInformation.ISP, string(geoString), geoInformation.Host, geoInformation.IP, geoInformation.RDNS, geoInformation.ASN, geoInformation.CountryName, geoInformation.CountryCode, geoInformation.RegionName, geoInformation.RegionCode, geoInformation.City, geoInformation.PostalCode, geoInformation.ContinentName, geoInformation.ContinentCode, geoInformation.Latitude, geoInformation.Longitude, geoInformation.Metro_code, geoInformation.Timezone, geoInformation.Datetime)
				finalString := " ON CONFLICT (address) DO UPDATE SET address = excluded.address, isp = excluded.isp, geo_data = excluded.geo_data, host = excluded.host, ip = excluded.ip, rdns = excluded.rdns, asn = excluded.asn, country_name = excluded.country_name, country_code = excluded.country_code, region_name = excluded.region_name, region_code = excluded.region_code, city = excluded.city, postal_code = excluded.postal_code, continent_name = excluded.continent_name, continent_code = excluded.continent_code, latitude = excluded.latitude, longitude = excluded.longitude, metro_code = excluded.metro_code, timezone = excluded.timezone, datetime = excluded.datetime"
				totalString := initialString + queryString + finalString

				_, err = database.DB.Exec(totalString)
				if err != nil {
					log.Printf("[ERROR] adding new validator geodata: %v", err)
				}

				log.Printf("Validator %v / %v - ISP: %v", address.String, validatorIP, geoInformation.ISP)

			}

			time.Sleep(time.Second * 5)
		}

	}

}

func getGeoData(ip string) ValidatorGeoData {

	url := "https://tools.keycdn.com/geo.json?host=" + ip

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "keycdn-tools:https://tools.keycdn.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}

	var geo ValidatorGeoData

	geoDataString := gjson.Get(string(body), "data.geo")
	if geoDataString.Exists() {
		json.Unmarshal([]byte(string(geoDataString.String())), &geo)
	}

	return geo
}

func getValidatorIP(input string) string {

	s := strings.Split(input, "/")

	if len(s) > 2 {
		return s[2]
	}

	return ""
}
