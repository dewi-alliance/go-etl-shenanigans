package scheduler

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"dewietl/database"

	"github.com/robfig/cron/v3"
)

// Create new cron
var scheduleCron = cron.New()

func Start() {

	REBUILD := flag.Bool("rebuild", false, "Rebuild the hotspot cache")

	flag.Parse()

	if *REBUILD == true {
		log.Println("Rebuilding cache!")
		go buildRewardCache()
	}

	// Slowly pulls geo data from helium api and inserts into the database
	go updateGetLocations()
	go updateValidatorGeoData()
	go updateMakersData()

}

func runScheduler(c *cron.Cron) {

	// Run yesterday job every day at midnight
	_, err := scheduleCron.AddFunc("0 0 * * *", func() {
		buildYesterdayCache()
	})

	if err != nil {
		log.Println(err)
	}

	// Start all cronjobs
	scheduleCron.Start()
}

// Helper functions

// getAllHotspots returns a slice of all hotspots in the database
func getAllHotspots() map[string]string {

	var address, location sql.NullString
	hotspots := make(map[string]string)

	rows, err := database.DB.Query("SELECT address, location FROM gateway_inventory ORDER BY first_block DESC")
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&address, &location)
		if err != nil {
			log.Printf("[ERROR] getAllHotspots(): %v", err)
		}

		if address.Valid == true && location.Valid == true {
			hotspots[address.String] = location.String
		}
	}

	return hotspots
}

// checkIfLocationExistsInDatabase checks if a location exists in the database
func checkIfLocationExistsInDatabase(location string) bool {

	rows, err := database.DB.Query("SELECT location FROM locations WHERE location = $1", location)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	defer rows.Close()

	exists := false

	for rows.Next() {
		exists = true
	}

	return exists
}

// geHotspotGeoLocation returns the geo data of a hotspot
func geHotspotGeoLocation(address string) GeoResponse {

	response, err := http.Get("https://api.helium.io/v1/hotspots/" + address)

	if err != nil {
		fmt.Println(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("[ERROR] geHotspotGeoLocation(): %v", err)
	}

	var res GeoResponse
	json.Unmarshal(responseData, &res)

	return res
}
