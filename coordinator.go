package main

import (
	"strconv"
)

var channelMap map[HashKey]chan string

func initGlobals() {
	channelMap = make(map[HashKey]chan string)
}

func coordinator() {
	defer wg.Done()
	initGlobals()

	for i := 0; i < len(nodeList); i++ {
		key := nodeList[i]
		channelMap[key] = make(chan string)
		wg.Add(1)
		go nodeWorker(key)
	}

	for elem := range coordinateChan {
		if elem == 20 {
			closeAllChannels()
		} else {
			channelMap[nodeList[elem]] <- strconv.Itoa(elem)
		}
	}
}

func closeAllChannels() {
	ticker.Stop()
	for _, channel := range channelMap {
		close(channel)
	}
	close(coordinateChan)
}
