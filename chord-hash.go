package main

import (
	"crypto/md5"
	"fmt"
	"hash/fnv"
	"math/rand"
	"sort"
	"time"
)

// HashKey - Holds the hash value of type uint32
type HashKey uint32

//HashKeyOrder - To sort hashkey
type HashKeyOrder []HashKey

func (h HashKeyOrder) Len() int           { return len(h) }
func (h HashKeyOrder) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h HashKeyOrder) Less(i, j int) bool { return h[i] < h[j] }

var nodeList []HashKey

func generateRandomID(size int) {

	for i := 0; i < size; i++ {
		key := genKey(randString())
		nodeList = append(nodeList, key)
	}
	sort.Sort(HashKeyOrder(nodeList))
	fmt.Println(nodeList)
}

var alphabets = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func randString() string {
	b := make([]rune, 15)
	for i := range b {
		rand.Seed(time.Now().UTC().UnixNano())
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
