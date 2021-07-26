package main

import (
	"dewietl/database"
	"dewietl/scheduler"
	"time"
)

func main() {

	// Start database
	database.Start()

	// Start scheduler
	scheduler.Start()

	for {
		time.Sleep(time.Second)
	}
}
