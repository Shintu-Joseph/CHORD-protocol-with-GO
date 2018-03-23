package main

import (
	"fmt"
	"strconv"
	"sync"
)

func coordinator(nodeList []uint64, wg *sync.WaitGroup, mainChannel chan string) {

	mainChannel <- "hi"
	defer wg.Done()

	fmt.Println("Inside coordinator")

	m := make(map[string]chan string)

	for i := 0; i < len(nodeList); i++ {

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			key := strconv.FormatUint(nodeList[i], 10)
			m[key] = make(chan string)
			m[key] <- "sssss"

			////range here///

			// message := <-c1 + strconv.Itoa(int(i))
			// r2, _ := json.Marshal(message)
			// c2 <- r2
			// fmt.Printf("Received message %d", i)
		}(i)
	}

	select {
	default:
		fmt.Println(m)
	}

	// var messages []string
	// j := 0
	// for i := range c2 {
	// 	str := i
	// 	var jresp string
	// 	json.Unmarshal([]byte(str), &jresp)
	// 	messages[j] = jresp

	// 	j++
	// 	if j == N {
	// 		fmt.Println("All messages received::", messages)
	// 		wg.Done()
	// 	}
	// }

	fmt.Println("Last line of coordinator")
}
