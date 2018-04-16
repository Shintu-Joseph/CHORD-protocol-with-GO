package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

//var channelMap map[HashKey]chan string
var channelMap map[HashKey]chan string

func initGlobals() {
	//channelMap = make(map[HashKey]chan string)
	channelMap = make(map[HashKey]chan string)
}

func coordinator() {
	defer wg.Done()
	initGlobals()

	for i := 0; i < len(nodeList); i++ {
		key := nodeList[i]
		channelMap[key] = make(chan string, 5)
		wg.Add(1)
		go nodeWorker(key, true)
	}

	//Send message to sponsor
	for message := range coordinateChan {
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(message), &dat); err != nil {
			panic(err)
		}
		if dat["Do"] == "join-ring" {

			key := genKey(randString())

			fmt.Println("here")
			channelMap[key] = make(chan string, 5)
			wg.Add(1)
			go nodeWorker(key, false)
			fmt.Println("here1")
			channelMap[key] <- message
			channelMap[key] <- initRingFingMessage()
		}
		if dat["Do"] == "leave-ring" {
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "stabilize-ring" {
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "init-ring-fingers" {
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "fix-ring-fingers" {
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "ring-notify" {
			res := doRespondToMsgs{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "get-ring-fingers" {
			res := doRespondToMsgs{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "find-ring-successor" {
			res := findRingSPMsg{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "find-ring-predecessor" {
			res := findRingSPMsg{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "put" {
			/*res := putMsg{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			*/
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "get" {
			/*res := getRemMsgs{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			*/
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
		}
		if dat["Do"] == "remove" {
			/*res := getRemMsgs{}
			json.Unmarshal([]byte(message), &res)
			sponsor := res.RespondTO
			*/
			sponsor := nodeList[rand.Intn(len(nodeList))]
			channelMap[sponsor] <- message
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
