package scheduler

import (
	"database/sql"
	"dewietl/database"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
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

				// Sort days
				keys := make([]string, 0, len(hotspotRewards))
				for k := range hotspotRewards {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				// // Fill the map

				i = 0
				for _, k := range keys {
					// for day, reward := range hotspotRewards {
					i++
					total--
					dayString := k[:10]

					// avoid adding current day
					if dayString != todayDate {

						dayParsed, _ := time.Parse("2006-01-02", dayString)
						dayTimestamp := dayParsed.Unix()

						rewardString := strconv.Itoa(hotspotRewards[k])
						dayParsedString := strconv.Itoa(int(dayTimestamp))

						query := `INSERT INTO rewards_by_day (address, date, amount) VALUES ('` + h + `', '` + dayParsedString + `', '` + rewardString + `')`

						_, err := database.DB.Exec(query)
						if err != nil {
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

		day, reward := getYesterdayRewards(h)

		query := `INSERT INTO rewards_by_day (address, date, amount) VALUES ('` + h + `', '` + day + `', '` + reward + `')`
		_, err := database.DB.Exec(query)
		if err != nil {
			log.Printf("[ERROR] error when inserting daily rewards: %v", err)
		}

		log.Printf("%v/%v Added yesterday rewards for %v total -> %v", currentHotspot, totalHotspots, h, reward)
	}

}

func getHostpotRewards(hash string) map[string]int {

	type Reward struct {
		Day    string `json:"day"`
		Reward int    `json:"reward"`
	}

	rewardList := make([]Reward, 0)

	var amount, timestamp sql.NullInt64
	rewards := make(map[string]int, 0)
	rewardsSorted := make(map[string]int, 0)
	finalList := make(map[string]int, 0)

	rows, err := database.DB.Query("SELECT amount, time FROM rewards WHERE gateway = $1 ORDER BY time", hash)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	defer rows.Close()

	j := 0
	for rows.Next() {
		j++
		err := rows.Scan(&amount, &timestamp)
		if err != nil {
			log.Printf("[ERROR] getAllHotspotsAssertions(): %v", err)
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
	rows.Close()

	if j > 0 {

		// Sort days
		keys := make([]string, 0, len(rewards))
		for k := range rewards {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Fill the map
		for _, k := range keys {
			day := strings.TrimSpace(k)
			rewardsSorted[day] = rewards[k]
			rewardList = append(rewardList, Reward{day, rewards[k]})
		}

		const layout = "2006-01-02"

		firstDay := ""
		for i, value := range rewardList {
			if i == 0 {
				firstDay = value.Day
			}
		}

		// Build all days
		firstDayTM, _ := time.Parse(layout, firstDay)
		lastDayTM := time.Now().AddDate(0, 0, -1)
		totalDays := lastDayTM.Sub(firstDayTM).Hours() / 24

		for i := 0; i <= int(totalDays); i++ {
			day := firstDayTM.AddDate(0, 0, i).Format(layout)
			finalList[day] = rewardsSorted[day]

		}

	}

	return finalList
}

func getYesterdayRewards(hash string) (string, string) {

	// Get beginning of the day timestamp
	year, month, day := time.Now().Date()
	beginningOfDayTimestamp := time.Date(year, month, day-1, 0, 0, 0, 0, time.UTC).Unix()
	endOfDayTimestamp := time.Date(year, month, day-1, 23, 59, 59, 0, time.UTC).Unix()

	var amount sql.NullInt64

	row := database.DB.QueryRow("SELECT sum(amount) FROM rewards WHERE gateway = $1 AND time >= $2 AND time <= $3", hash, beginningOfDayTimestamp, endOfDayTimestamp)
	err := row.Scan(&amount)

	if err != nil {
		log.Printf("[ERROR]d error: %v ", err)
	}

	dayString := fmt.Sprintf("%v", beginningOfDayTimestamp)
	rewards := fmt.Sprintf("%v", amount.Int64)
	return dayString, rewards
}
