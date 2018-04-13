package main

import (
	"fmt"
	"math"
)

func nodeWorker(key HashKey) {
	defer wg.Done()

	nodeChan := channelMap[key]
	bucket := make(map[HashKey]string)
	fingerTable := make([]HashKey, 32)
	initialRingSimulator(fingerTable, key)
	for elem := range nodeChan {
		bucket[key] = elem
		if elem == "15" {
			fmt.Println(elem, fingerTable)
			channelMap[HashKey(916619801)] <- "in 4"
		}
	}
}

func getSuccessor(key HashKey) {

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
