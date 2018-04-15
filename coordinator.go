package main

import (
	"encoding/json"
)

//var channelMap map[HashKey]chan string
var channelMap map[HashKey]chan []byte

func initGlobals() {
	//channelMap = make(map[HashKey]chan string)
	channelMap = make(map[HashKey]chan []byte)
}

func checkKey(key HashKey) HashKey {
	key = genKey(randString())
	for _, node := range nodeList {
		if node == key {
			checkKey(key)
		}
	}
	return key
}

func coordinator() {
	defer wg.Done()
	initGlobals()

	for i := 0; i < len(nodeList); i++ {
		key := nodeList[i]
		//channelMap[key] = make(chan string)
		channelMap[key] = make(chan []byte)
		wg.Add(1)
		go nodeWorker(key)
	}

	/*for elem := range coordinateChan {
		if elem == 20 {
			closeAllChannels()
		} else {
			channelMap[nodeList[elem]] <- strconv.Itoa(elem)
		}
	}*/

	//Send message to sponsor
	for message := range coordinateChan {
		var dat map[string]interface{}
		if err := json.Unmarshal(message, &dat); err != nil {
			panic(err)
		}
		if dat["Do"] == "join-ring" {
			key := genKey(randString())

			for _, node := range nodeList {
				if node == key {
					key = checkKey(key)
				}
			}
			wg.Add(1)
			go nodeWorker(key)
			channelMap[key] <- message
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
