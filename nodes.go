package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
)

func nodeWorker(key HashKey, buildRing bool) {
	defer wg.Done()

	nodeChan := channelMap[key]
	bucket := make(map[HashKey]string)
	fingerTable := make([]HashKey, 32)

	recipient := key

	if buildRing {
		initialRingSimulator(fingerTable, key)
	}
	successor := fingerTable[0]
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

			}
		case "get-ring-fingers":
			{
				getRingFingers(message, fingerTable)
			}
		case "put":
			{
				putData(message, key, bucket)
			}
		case "get":
			{
				getData(message, key, bucket)
			}
		case "remove":
			{
				removeData(message, key, bucket)
			}
		case "leave-ring":
			{

				res := leaveRingMsg{}
				json.Unmarshal([]byte(message), &res)
				if res.Mode == "orderly" {
					prepareToLeaveRing(successor, key, bucket)
				}
				leaveRing(key)
			}
		case "update-bucket":
			{
				bucket = updateBucket(message, bucket)
			}
		case "find-ring-predecessor":
			{

				res := findRingSPMsg{}
				json.Unmarshal([]byte(message), &res)

				n := res.RespondTO
				ID := res.TargetID
				predecessorToRespond := getPredecessor(n, ID, fingerTable)
				channelMap[n] <- strconv.FormatUint(uint64(predecessorToRespond), 10)
			}

		}

	}
}

func updateBucket(message string, bucket map[HashKey]string) map[HashKey]string {

	res := updateBucketsMsg{}
	json.Unmarshal([]byte(message), &res)
	for k, v := range res.BucketData {
		bucket[k] = v
	}
	return bucket
}

func leaveRing(key HashKey) {
	wg.Done()
	close(channelMap[key])
	delete(channelMap, key)
	removeNodeFromList(key)

}

func prepareToLeaveRing(successor HashKey, node HashKey, bucket map[HashKey]string) {
	channelMap[successor] <- updateBucketMessage(bucket)
	sponsor := nodeList[rand.Intn(len(nodeList))]

	channelMap[sponsor] <- triggerPredecessorMessage(sponsor, node)
	predecessorBytes := <-channelMap[sponsor]
	fin, _ := strconv.ParseUint(predecessorBytes, 10, 64)
	predecessor := HashKey(uint32(fin))

	fmt.Println(predecessor, sponsor, node)
	fmt.Println(nodeList)
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

	}

	closestNode := findNearestPreceedingNode(recipient, fingerTable)
	channelMap[closestNode] <- triggerSuccesorMessage(closestNode, recipient)

	successorBytes := <-channelMap[closestNode]
	fin, _ := strconv.ParseUint(successorBytes, 10, 64)
	successor := HashKey(uint32(fin))

	return successor

}

func getPredecessor(sponsor HashKey, recipient HashKey, fingerTable []HashKey) HashKey {
	if recipient > sponsor && recipient <= fingerTable[0] {
		return sponsor

	} else {
		closestNode := findNearestPreceedingNode(recipient, fingerTable)
		if closestNode == recipient {
			return closestNode
		}
		channelMap[closestNode] <- triggerPredecessorMessage(closestNode, recipient)

		predecessorBytes := <-channelMap[closestNode]
		fin, _ := strconv.ParseUint(predecessorBytes, 10, 64)
		predecessor := HashKey(uint32(fin))

		return predecessor
	}
}

func putData(message string, key HashKey, bucket map[HashKey]string) {
	res := putMsg{}
	json.Unmarshal([]byte(message), &res)
	n := res.RespondTO
	ID := res.Data.Key
	channelMap[n] <- triggerSuccesorMessage(n, ID)
	successorBytes := <-channelMap[n]
	fin, _ := strconv.ParseUint(successorBytes, 10, 64)
	successor := HashKey(uint32(fin))
	if key == successor {
		bucket[ID] = res.Data.Value
		fmt.Println("bucket inserted: ", bucket[ID])
	} else {
		channelMap[successor] <- message
	}
}

func getData(message string, key HashKey, bucket map[HashKey]string) {
	res := getRemMsgs{}
	json.Unmarshal([]byte(message), &res)
	n := res.RespondTO
	ID := res.Data.Key
	channelMap[n] <- triggerSuccesorMessage(n, ID)
	successorBytes := <-channelMap[n]
	fin, _ := strconv.ParseUint(successorBytes, 10, 64)
	successor := HashKey(uint32(fin))
	if key == successor {
		fmt.Println("successor = ", successor)
		fmt.Println("key =", key)
		value := bucket[ID]
		fmt.Println("Value = ", value)
	} else {
		channelMap[successor] <- message
	}
}

func removeData(message string, key HashKey, bucket map[HashKey]string) {
	res := getRemMsgs{}
	json.Unmarshal([]byte(message), &res)
	n := res.RespondTO
	ID := res.Data.Key
	channelMap[n] <- triggerSuccesorMessage(n, ID)
	successorBytes := <-channelMap[n]
	fin, _ := strconv.ParseUint(successorBytes, 10, 64)
	successor := HashKey(uint32(fin))
	if key == successor {
		//fmt.Println("Deleted")
		delete(bucket, ID)
		fmt.Println("Deleted")
	} else {
		channelMap[successor] <- message
	}
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
		if node < key {
			return node
		}
	}
	return tempTable[0]
}
