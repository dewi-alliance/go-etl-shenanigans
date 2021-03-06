package handler

import (
	"database/sql"
	"dewietl/database"
	"log"

	"github.com/labstack/echo"
)

func GetLocationAddress(c echo.Context) error {

	// location, long_street, short_street, long_city, short_city, long_state, short_state, long_country, short_country, search_city, city_id
	type Address struct {
		Location     string `json:"location"`
		LongStreet   string `json:"long_street"`
		ShortStreet  string `json:"short_street"`
		LongCity     string `json:"long_city"`
		ShortCity    string `json:"short_city"`
		LongState    string `json:"long_state"`
		ShortState   string `json:"short_state"`
		LongCountry  string `json:"long_country"`
		ShortCountry string `json:"short_country"`
		SearchCity   string `json:"search_city"`
		CityId       string `json:"city_id"`
	}

	hash := c.Param("hash")

	var location, long_street, short_street, long_city, short_city, long_state, short_state, long_country, short_country, search_city, city_id sql.NullString

	err := database.DB.QueryRow(`SELECT location, long_street, short_street, long_city, short_city, long_state, short_state, long_country, short_country, search_city, city_id FROM locations WHERE location = $1`, hash).Scan(&location, &long_street, &short_street, &long_city, &short_city, &long_state, &short_state, &long_country, &short_country, &search_city, &city_id)
	if err != nil {
		log.Printf("[ERROR] GetLocationAddress(): %v", err)
	}

	address := Address{location.String, long_street.String, short_street.String, long_city.String, short_city.String, long_state.String, short_state.String, long_country.String, short_country.String, search_city.String, city_id.String}

	return c.JSON(200, address)
}

func GetLocationHotspot(c echo.Context) error {

	// location, long_street, short_street, long_city, short_city, long_state, short_state, long_country, short_country, search_city, city_id
	type Address struct {
		Location     string `json:"location"`
		LongStreet   string `json:"long_street"`
		ShortStreet  string `json:"short_street"`
		LongCity     string `json:"long_city"`
		ShortCity    string `json:"short_city"`
		LongState    string `json:"long_state"`
		ShortState   string `json:"short_state"`
		LongCountry  string `json:"long_country"`
		ShortCountry string `json:"short_country"`
		SearchCity   string `json:"search_city"`
		CityId       string `json:"city_id"`
	}

	hash := c.Param("hash")

	var location, long_street, short_street, long_city, short_city, long_state, short_state, long_country, short_country, search_city, city_id sql.NullString

	query := `SELECT
		lo.location, 
		long_street, 
		short_street, 
		long_city, 
		short_city, 
		long_state, 
		short_state, 
		long_country, 
		short_country, 
		search_city, 
		city_id
		FROM
		locations lo
		INNER JOIN gateway_inventory gi 
		ON gi.location = lo.location
		WHERE  gi.address = '` + hash + `'`

	err := database.DB.QueryRow(query).Scan(&location, &long_street, &short_street, &long_city, &short_city, &long_state, &short_state, &long_country, &short_country, &search_city, &city_id)
	if err != nil {
		// log.Printf("[ERROR] GetLocationHotspot(): %v - %v", err, hash)
	}

	address := Address{location.String, long_street.String, short_street.String, long_city.String, short_city.String, long_state.String, short_state.String, long_country.String, short_country.String, search_city.String, city_id.String}

	return c.JSON(200, address)
}
