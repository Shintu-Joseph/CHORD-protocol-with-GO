package main

import (
	"fmt"
	"hash/crc32"
	"math/rand"
)

func hashID(ringOrder int) uint64 {
	generateRandomID(40)
	return 0
}

func generateRandomID(size int) {
	nodeList := [20]uint32{}
	for i := 0; i < len(nodeList); i++ {
		nodeList[i] = hasKey(randString())
	}
	fmt.Println(nodeList)
}

var alphabets = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString() string {
	b := make([]rune, 15)
	for i := range b {
		b[i] = alphabets[rand.Intn(52)]
	}
	return string(b)
}

func hasKey(obj string) uint32 {
	var scratch [64]byte
	if len(obj) < 64 {

		copy(scratch[:], obj)
	}
	return crc32.ChecksumIEEE(scratch[:len(obj)])
}
