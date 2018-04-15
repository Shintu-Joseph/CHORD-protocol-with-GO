package main

import (
	"encoding/json"
	"fmt"
)

//var channelMap map[HashKey]chan string
var channelMap map[HashKey]chan []byte

func initGlobals() {
	//channelMap = make(map[HashKey]chan string)
	channelMap = make(map[HashKey]chan []byte)
}

func coordinator() {
	defer wg.Done()
	initGlobals()

	for i := 0; i < len(nodeList); i++ {
		key := nodeList[i]
		//channelMap[key] = make(chan string)
		channelMap[key] = make(chan []byte)
		wg.Add(1)
		go nodeWorker(key, true)
	}

	// for elem := range coordinateChan {
	// 	if elem == 20 {
	// 		closeAllChannels()
	// 	} else {
	// 		channelMap[nodeList[elem]] <- strconv.Itoa(elem)
	// 	}
	// }

	//Send message to sponsor
	for message := range coordinateChan {
		var dat map[string]interface{}
		if err := json.Unmarshal(message, &dat); err != nil {
			panic(err)
		}
		if dat["Do"] == "join-ring" {

			key := genKey(randString())

			fmt.Println("here")
			channelMap[key] = make(chan []byte)
			wg.Add(1)
			go nodeWorker(key, false)
			fmt.Println("here1")
			channelMap[key] <- message
		}
	}
}

func closeAllChannels() {
	fmt.Println("s")
	ticker.Stop()
	for _, channel := range channelMap {
		close(channel)
	}
	close(coordinateChan)
}
