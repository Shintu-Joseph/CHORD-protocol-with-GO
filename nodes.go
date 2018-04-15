package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

var count int

func nodeWorker(key HashKey) {
	defer wg.Done()

	nodeChan := channelMap[key]
	bucket := make(map[HashKey]string)
	fingerTable := make([]HashKey, 32)
	recipient := key
	if count == 0 {
		initialRingSimulator(fingerTable, key)
		count = count + 1
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
				if res.Do == "join-ring" {
					sponsorKey := res.Sponsor
					joinRing(sponsorKey, recipient)
				}
			}
		case "find-ring-successor":
			{
				res := findRingSPMsg{}
				json.Unmarshal(message, &res)
				if res.Do == "find-ring-successor" {
					n := res.RespondTO
					ID := res.TargetID
					successor := getSuccessor(n, ID, fingerTable)
				}

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
	findSuccesorM := &findRingSPMsg{
		Do:        "find-ring-successor",
		RespondTO: sponsor,
		TargetID:  recipient,
	}
	fsMessage, _ := json.Marshal(findSuccesorM)
	channelMap[sponsor] <- fsMessage
	successor := getSuccessor(sponsor, recipient)
	//init ring fingers
	initRingFingers(recipient, successor)
	//append to nodeList
	joinChord(recipient)

}

//init ring fingers of joining node
func initRingFingers(recipient HashKey, successor HashKey) {
	//copy from successor
	fingerTable_recipient := fingerTable
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
			return successor
		}
	}
}

func getPredecessor(key HashKey) {

}

func initialRingSimulator(fingerTable []HashKey, key HashKey) {

	for i := 0; i < 32; i++ {
		fingerTable[i] = findNearestNode(HashKey((int(key) + int(math.Pow(2, float64(i)))) % int(math.Pow(2, 32))))
	}
}

func findNearestNode(key HashKey) HashKey {
	for _, node := range nodeList {
		if node >= key {
			fmt.Println(key, node)
			return node
		}
	}

	return nodeList[0]
}
