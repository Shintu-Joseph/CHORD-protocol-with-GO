package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
)

func nodeWorker(key HashKey, buildRing bool) {
	defer wg.Done()

	nodeChan := channelMap[key]
	//bucket := make(map[HashKey]string)
	fingerTable := make([]HashKey, 32)
	recipient := key
	successor := HashKey(0)
	if buildRing {
		initialRingSimulator(fingerTable, key)
	}

	for message := range nodeChan {

		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(message), &dat); err != nil {
			panic(err)
		}
		choice := dat["Do"]
		switch choice {
		case "join-ring":
			{
				res := joinRingMsg{}
				json.Unmarshal([]byte(message), &res)
				sponsorKey := res.Sponsor
				successor = joinRing(sponsorKey, recipient, fingerTable)

			}
		case "find-ring-successor":
			{

				res := findRingSPMsg{}
				json.Unmarshal([]byte(message), &res)

				n := res.RespondTO
				ID := res.TargetID
				successorToRespond := getSuccessor(n, ID, fingerTable)
				channelMap[n] <- strconv.FormatUint(uint64(successorToRespond), 10)

			}
		case "init-ring-fingers":
			{
				initRingFingers(key, successor, fingerTable)
				fmt.Println(fingerTable, key, successor)

			}
		case "get-ring-fingers":
			{
				getRingFingers(message, fingerTable)
			}

		}

	}
}

//coordinator instructs recipient node to join ring
func joinRing(sponsor HashKey, recipient HashKey, fingerTable []HashKey) HashKey {

	//find node successor

	channelMap[sponsor] <- triggerSuccesorMessage(sponsor, recipient)

	successorBytes := <-channelMap[sponsor]
	fin, _ := strconv.ParseUint(successorBytes, 10, 64)
	successor := HashKey(uint32(fin))

	//init ring fingers

	joinChord(recipient)

	return successor

}

func getRingFingers(message string, fingerTable []HashKey) {
	res := doRespondToMsgs{}
	json.Unmarshal([]byte(message), &res)

	recipient := res.RespondTO
	marshalledFing, _ := json.Marshal(fingerTable)
	channelMap[recipient] <- string(marshalledFing)
}

//init ring fingers of joining node
func initRingFingers(recipient HashKey, successor HashKey, fingerTable []HashKey) {

	channelMap[successor] <- getRingFingMessage(recipient)

	tempFingTable := []HashKey{}
	json.Unmarshal([]byte(<-channelMap[recipient]), &tempFingTable)

	//Update the finger table

	copyFingerTable(recipient, successor, fingerTable, tempFingTable)

}

//Add new node to nodeList
func joinChord(key HashKey) {
	nodeList = append(nodeList, key)
	sort.Sort(HashKeyOrder(nodeList))
}

func getSuccessor(sponsor HashKey, recipient HashKey, fingerTable []HashKey) HashKey {
	if recipient > sponsor && recipient < fingerTable[0] {
		return fingerTable[0]

	} else {
		closestNode := findNearestPreceedingNode(recipient, fingerTable)
		channelMap[closestNode] <- triggerSuccesorMessage(closestNode, recipient)

		successorBytes := <-channelMap[closestNode]
		fin, _ := strconv.ParseUint(successorBytes, 10, 64)
		successor := HashKey(uint32(fin))

		return successor
	}
}

func getPredecessor(key HashKey) {

}

func initialRingSimulator(fingerTable []HashKey, key HashKey) {

	for i := 0; i < 32; i++ {
		key := HashKey((int(key) + int(math.Pow(2, float64(i)))) % int(math.Pow(2, 32)))
		fingerTable[i] = findNearestSuccessorNode(key)
	}
}

func findNearestSuccessorNode(key HashKey) HashKey {
	for _, node := range nodeList {
		if node >= key {
			return node
		}
	}
	return nodeList[0]
}

func findNearestPreceedingNode(key HashKey, fingerTable []HashKey) HashKey {
	tempTable := fingerTable
	sort.Sort(sort.Reverse(HashKeyOrder(tempTable)))
	for _, node := range tempTable {
		if node <= key {
			return node
		}
	}
	return tempTable[0]
}
