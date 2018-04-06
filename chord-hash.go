package main

import (
	"crypto/md5"
	"fmt"
	"hash/fnv"
	"math/rand"
)

type HashKey uint32

func hashID(ringOrder int) uint64 {
	generateRandomID(40)
	return 0
}

func generateRandomID(size int) []HashKey {
	var nodeList []HashKey
	for i := 0; i < size; i++ {
		key := genKey(randString())
		nodeList = append(nodeList, key)
		fmt.Println(key)
	}

	return nodeList
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
	h := fnv.New32a()
	h.Write([]byte(obj))
	return h.Sum32()
}

func genKey(key string) HashKey {
	bKey := hashDigest(key)
	return hashVal(bKey[0:4])
}

func hashDigest(key string) [md5.Size]byte {
	return md5.Sum([]byte(key))
}

func hashVal(bKey []byte) HashKey {
	return ((HashKey(bKey[3]) << 24) |
		(HashKey(bKey[2]) << 16) |
		(HashKey(bKey[1]) << 8) |
		(HashKey(bKey[0])))
}
