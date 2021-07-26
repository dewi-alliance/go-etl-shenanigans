package scheduler

import (
	"fmt"
	"log"
	"strings"
	"time"

	"dewietl/database"
)

type GeoResponse struct {
	Data GeoData `json:"data"`
}

type GeoData struct {
	Address  string  `json:"id"`
	Name     string  `json:"name"`
	Lat      string  `json:"lat"`
	Lng      string  `json:"lng"`
	Location string  `json:"location"`
	Geo      GeoCode `json:"geocode"`
}

type GeoCode struct {
	ShortStreet  string `json:"short_street"`
	ShortState   string `json:"short_state"`
	ShortCountry string `json:"short_country"`
	ShortCity    string `json:"short_city"`
	LongStreet   string `json:"long_street"`
	LongState    string `json:"long_state"`
	LongCountry  string `json:"long_country"`
	LongCity     string `json:"long_city"`
	CityID       string `json:"city_id"`
}

func updateGetLocations() {
	for {
		hotspots := getAllHotspots()

		for address := range hotspots {

			locationData := geHotspotGeoLocation(address)

			location := strings.Replace(locationData.Data.Location, "'", "''", -1)
			shortStreet := strings.Replace(locationData.Data.Geo.ShortStreet, "'", "''", -1)
			shortState := strings.Replace(locationData.Data.Geo.ShortState, "'", "''", -1)
			shortCountry := strings.Replace(locationData.Data.Geo.ShortCountry, "'", "''", -1)
			shortCity := strings.Replace(locationData.Data.Geo.ShortCity, "'", "''", -1)
			longStreet := strings.Replace(locationData.Data.Geo.LongStreet, "'", "''", -1)
			longState := strings.Replace(locationData.Data.Geo.LongState, "'", "''", -1)
			longCountry := strings.Replace(locationData.Data.Geo.LongCountry, "'", "''", -1)
			longCity := strings.Replace(locationData.Data.Geo.LongCity, "'", "''", -1)
			cityID := strings.Replace(locationData.Data.Geo.CityID, "'", "''", -1)

			initialString := "INSERT INTO locations (location, short_street, short_state, short_country, short_city, long_street, long_state, long_country, long_city, city_id) VALUES "
			queryString := fmt.Sprintf(" ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v') ", location, shortStreet, shortState, shortCountry, shortCity, longStreet, longState, longCountry, longCity, cityID)
			finalString := " ON CONFLICT (location) DO UPDATE SET short_street = excluded.short_street, short_state = excluded.short_state, short_country = excluded.short_country, short_city = excluded.short_city, long_street = excluded.long_street, long_state = excluded.long_state, long_country = excluded.long_country, long_city = excluded.long_city, city_id = excluded.city_id"
			totalString := initialString + queryString + finalString

			_, err := database.DB.Exec(totalString)
			if err != nil {
				log.Printf("[ERROR] updateGetLocations(): %v - %v - %v - %v - %v - %v - %v - %v - %v - %v - %v -> %v ", err, location, shortStreet, shortState, shortCountry, shortCity, longStreet, longState, longCountry, longCity, cityID, totalString)
			}

			log.Printf("%v - %v - %v", longCountry, longCity, location)

			time.Sleep(time.Second)
		}

	}
}
