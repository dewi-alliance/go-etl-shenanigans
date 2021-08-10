package scheduler

import (
	"database/sql"
	"dewietl/database"
	"fmt"
	"log"
	"strconv"
	"time"
)

// buildRewardCache created the cache for "full-days", every time this functions
func buildRewardCache() {

	hotspots := getAllHotspots()
	totalHotspots := len(hotspots)
	currentHotspot := 0

	for h := range hotspots {

		i := 0
		currentHotspot++

		hotspotRewards := getHostpotRewards(h)

		// Only insert if hotspot has rewards
		// Do not save the current day

		todayDate := fmt.Sprintf("%v", time.Now().Format("2006-01-02"))

		// if hotspot has rewards
		if len(hotspotRewards) > 0 {

			// total amount of rewards for this hotspor
			total := len(hotspotRewards)

			// disable hotspots with only today as date
			if total > 1 {

				i = 0
				for day, reward := range hotspotRewards {
					i++
					total--
					dayString := day[:10]

					// avoid adding current day
					if dayString != todayDate {

						dayParsed, _ := time.Parse("2006-01-02", dayString)
						dayTimestamp := dayParsed.Unix()

						rewardString := strconv.Itoa(reward)
						dayParsedString := strconv.Itoa(int(dayTimestamp))

						query := `INSERT INTO rewards_by_day (address, date, amount) 
								  SELECT '` + h + `', '` + dayParsedString + `', ` + rewardString + `
								  WHERE NOT EXISTS (
								  SELECT 1 FROM rewards_by_day WHERE address='` + h + `' AND date='` + dayParsedString + `' AND amount=` + rewardString + `);`

						_, err := database.DB.Exec(query)
						if err != nil {
							// log.Printf("\n\n\n %v", totalQuery)
							log.Printf("[ERROR] error when inserting daily rewards: %v", err)
						}
					}
				}

			}

			log.Printf("%v/%v [%v] Added daily rewards for %v", currentHotspot, totalHotspots, i, h)

		}
	}
}

// buildYesterdayCache builds the reward cache only for yesterday's date and should run everyday at midght
// potentially would be better to do only one query and parse everything on go instead of on query per hotspot
func buildYesterdayCache() {

	hotspots := getAllHotspots()
	totalHotspots := len(hotspots)
	currentHotspot := 0

	for h := range hotspots {

		currentHotspot++

		hotspotRewards := getYesterdayRewards(h)

		// if hotspot has rewards
		if len(hotspotRewards) > 0 {

			for day, reward := range hotspotRewards {

				rewardString := strconv.Itoa(reward)

				query := `INSERT INTO rewards_by_day (address, date, amount) 
				SELECT '` + h + `', '` + day + `', ` + rewardString + `
				WHERE NOT EXISTS (
				SELECT 1 FROM rewards_by_day WHERE address='` + h + `' AND date='` + day + `' AND amount=` + rewardString + `);`

				_, err := database.DB.Exec(query)
				if err != nil {
					// log.Printf("\n\n\n %v", totalQuery)
					log.Printf("[ERROR] error when inserting daily rewards: %v", err)
				}

			}
		}

		log.Printf("%v/%v Added yesterday rewards for %v", currentHotspot, totalHotspots, h)
		time.Sleep(time.Millisecond * 10)
	}

}

func getYesterdayRewards(hash string) map[string]int {

	// Get beginning of the day timestamp
	year, month, day := time.Now().Date()
	beginningOfDayTimestamp := time.Date(year, month, day-1, 0, 0, 0, 0, time.UTC).Unix()
	endOfDayTimestamp := time.Date(year, month, day-1, 23, 59, 59, 0, time.UTC).Unix()

	var amount sql.NullInt64
	rewards := make(map[string]int, 0)

	row := database.DB.QueryRow("SELECT sum(amount) FROM rewards WHERE gateway = $1 AND time >= $2 AND time <= $3", hash, beginningOfDayTimestamp, endOfDayTimestamp)
	err := row.Scan(&amount)

	if err != nil {
		log.Printf("[ERROR]d error: %v ", err)
	}

	dayString := fmt.Sprintf("%v", beginningOfDayTimestamp)
	rewards[dayString] = int(amount.Int64)
	return rewards
}

func getHostpotRewards(hash string) map[string]int {

	var amount, timestamp sql.NullInt64
	rewards := make(map[string]int, 0)

	rows, err := database.DB.Query("SELECT amount, time FROM rewards WHERE gateway = $1 ORDER BY time", hash)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&amount, &timestamp)
		if err != nil {
			log.Printf("[ERROR]: %v", err)
		}

		// get the date
		tm := time.Unix(timestamp.Int64, 0)
		dateTime := fmt.Sprintln(tm.Format("2006-01-02"))

		if _, ok := rewards[dateTime]; ok {
			rewards[dateTime] += int(amount.Int64)
		} else {
			rewards[dateTime] = int(amount.Int64)
		}
	}

	return rewards
}
