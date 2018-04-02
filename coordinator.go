package main

import (
	"fmt"
	"strconv"
)

var channelMap map[string]chan string

func coordinator(nodeList []uint64) {
	defer wg.Done()

	fmt.Println("Inside coordinator")
	channelMap = make(map[string]chan string)
	key2 := strconv.FormatUint(nodeList[2], 10)
	key0 := strconv.FormatUint(nodeList[0], 10)
	key3 := strconv.FormatUint(nodeList[3], 10)

	for i := 0; i < len(nodeList); i++ {
		key := strconv.FormatUint(nodeList[i], 10)
		channelMap[key] = make(chan string)
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			for elem := range channelMap[key] {
				fmt.Println(i, elem)
				if elem == "send coordin" {
					channelMap[key2] <- "send value to node 2"
					closeAllChannels()
				}
			}
		}(i)
	}

	channelMap[key0] <- "0"
	channelMap[key3] <- "send coordin"

	fmt.Println("Last line of coordinator")
}

func closeAllChannels() {
	for _, channel := range channelMap {
		close(channel)
	}
}
