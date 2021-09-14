package scheduler

import (
	"dewietl/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Maker struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type MakerData struct {
	Code    int32   `json:"code"`
	Data    []Maker `json:"data"`
	Success bool    `json:"success"`
}

func updateMakersData() {

	for {

		makers := getMakers()

		for _, maker := range makers.Data {

			initialString := "INSERT INTO makers (name, address) VALUES "
			queryString := fmt.Sprintf(" ('%v', '%v') ", maker.Name, maker.Address)
			finalString := " ON CONFLICT (address) DO UPDATE SET address = excluded.address"
			totalString := initialString + queryString + finalString

			_, err := database.DB.Exec(totalString)
			if err != nil {
				log.Printf("[ERROR] adding new validator geodata: %v", err)
			}

		}

		time.Sleep(time.Minute * 30)
	}

}

func getMakers() MakerData {

	url := "https://onboarding.dewi.org/api/v2/makers"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "dewi-shenanigans")

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

	var maker MakerData

	jsonErr := json.Unmarshal(body, &maker)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return maker
}
