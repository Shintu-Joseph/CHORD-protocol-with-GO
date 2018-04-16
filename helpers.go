package main

import (
	"fmt"
)

func copyFingerTable(node HashKey, successor HashKey, fingerTable []HashKey, successorFingerTable []HashKey) {
	fmt.Println(len(successorFingerTable))
	for i, v := range successorFingerTable {
		fingerTable[i] = v

	}
}
