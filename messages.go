package main

import (
	"math/rand"
	"time"
)

var ticker *time.Ticker

func injectRequests() {
	ticker = time.NewTicker(1500 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			coordinateChan <- generateMessages(t)
		}
	}()
}

func generateMessages(timeSeed time.Time) int {
	return randomGenerator(timeSeed, 0, 21)
}

func randomGenerator(timeSeed time.Time, min int, max int) int {
	rand.Seed(timeSeed.UTC().UnixNano())
	return rand.Intn(max-min) + min
}
