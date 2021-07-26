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

			initialString := "INSERT INTO validator_isp (address, isp, geo_data) VALUES "
			queryString := fmt.Sprintf(" ('%v', '%v', '%v') ", address, geoInformation.ISP, string(geoString))
			finalString := " ON CONFLICT (validator_isp) DO UPDATE SET address = excluded.address, isp = excluded.isp, geo_data = excluded.geo_data"
			totalString := initialString + queryString + finalString

			_, err = database.DB.Exec(totalString)
			if err != nil {
				log.Printf("[ERROR] adding new validator geo data()")
			}

			log.Println("Validator %v / %v - ISP: %v", address, validatorIP, geoInformation.ISP)

		}

		time.Sleep(time.Second * 5)
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
		log.Fatal(readErr)
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
