package main

func copyArray(fingerTable []HashKey, successorFingerTable []HashKey) {
	for i, v := range successorFingerTable {
		fingerTable[i] = v
	}
}
