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
	if buildRing {
		initialRingSimulator(fingerTable, key)
	}

	/*for elem := range nodeChan {
		bucket[key] = elem
		if elem == "15" {
			fmt.Println(elem, fingerTable)
			channelMap[HashKey(916619801)] <- "in 4"
		}
	}
	*/

	for message := range nodeChan {

		var dat map[string]interface{}
		if err := json.Unmarshal(message, &dat); err != nil {
			panic(err)
		}
		choice := dat["Do"]
		switch choice {
		case "join-ring":
			{
				res := joinRingMsg{}
				json.Unmarshal(message, &res)
				sponsorKey := res.Sponsor
				joinRing(sponsorKey, recipient)

			}
		case "find-ring-successor":
			{

				res := findRingSPMsg{}
				json.Unmarshal(message, &res)

				n := res.RespondTO
				ID := res.TargetID
				successor := getSuccessor(n, ID, fingerTable)
				channelMap[n] <- []byte(strconv.FormatUint(uint64(successor), 10))

			}
			/*case "leave-ring":
				return
			case "stabilize-ring":
				return
			*/
		}

	}
}

//coordinator instructs recipient node to join ring
func joinRing(sponsor HashKey, recipient HashKey) {
	//find node successor

	channelMap[sponsor] <- triggerSuccesorMessage(sponsor, recipient)

	successorBytes := <-channelMap[sponsor]
	fin, _ := strconv.ParseUint(string(successorBytes), 10, 64)
	successor := HashKey(uint32(fin))
	//init ring fingers
	initRingFingers(recipient, successor)
	//append to nodeList
	joinChord(recipient)

}

//init ring fingers of joining node
func initRingFingers(recipient HashKey, successor HashKey) {
	//copy from successor
	// fingerTable_recipient := fingerTable
}

//Add new node to nodeList
func joinChord(key HashKey) {
	nodeList = append(nodeList, key)
	sort.Sort(HashKeyOrder(nodeList))
}

/*func checkKey(key HashKey) HashKey {
	key = genKey(randString())
	for _, node := range nodeList {
		if node == key {
			checkKey(key)
		}
	}
	return key
}
*/

func getSuccessor(sponsor HashKey, recipient HashKey, fingerTable []HashKey) HashKey {
	for _, node := range fingerTable {
		if node > sponsor && recipient < node {
			successor := node
			fmt.Println(successor)
			return successor
		}
	}
	closestNode := findNearestPreceedingNode(recipient, fingerTable)
	channelMap[closestNode] <- triggerSuccesorMessage(closestNode, recipient)

	successorBytes := <-channelMap[closestNode]
	fin, _ := strconv.ParseUint(string(successorBytes), 10, 64)
	successor := HashKey(uint32(fin))

	return successor
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
