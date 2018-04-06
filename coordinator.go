package main

import (
	"fmt"
)

var channelMap map[HashKey]chan string

var bucketList map[string](map[string]string)

func coordinator(nodeList []HashKey) {
	defer wg.Done()

	fmt.Println("Inside coordinator")
	channelMap = make(map[HashKey]chan string)
	key2 := nodeList[2]
	key0 := nodeList[0]
	key3 := nodeList[3]

	for i := 0; i < len(nodeList); i++ {
		key := nodeList[i]
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
