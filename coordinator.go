package main

import (
	"fmt"
	"strconv"
)

var channelMap map[HashKey]chan string

var bucketList map[HashKey](map[HashKey]string)

var fingerTableList map[HashKey]HashKey

func initGlobals() {
	channelMap = make(map[HashKey]chan string)
	bucketList = make(map[HashKey](map[HashKey]string))
	fingerTableList = make(map[HashKey]HashKey)
}

func coordinator(nodeList []HashKey) {

	defer wg.Done()

	initGlobals()

	for i := 0; i < len(nodeList); i++ {
		key := nodeList[i]
		channelMap[key] = make(chan string)
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			initBucket(key)
			for elem := range channelMap[key] {
				fmt.Println(i, elem)
			}
		}(i)
	}

	for elem := range coordinateChan {

		if elem == 20 {
			closeAllChannels()
		} else {
			channelMap[nodeList[elem]] <- "send coordin" + strconv.Itoa(elem)

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

func initBucket(key HashKey) {
	bucketList[key] = make(map[HashKey]string)
}
