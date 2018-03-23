package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func coordinator(N int) {
	defer wg.Done()

	fmt.Println("Inside coordinator")

	//create coordinator channel
	c1 := make(chan string)
	c2 := make(chan []byte, N)
	c1 <- "Sending message"

	// m := make(map[string]chan string)
	// m["SS"] = make(chan string)
	// m["SS"] <- "sssss"

	for i := 0; i < N; i++ {
		//create random channel variable for chord nodes

		//m := make(map[string]chan string)
		//nodelist := generateRandomID(40)
		//m[strconv.Itoa(int(nodelist[i]))] = c2
		//m["k1"].(chan string) = make(chan string)
		//m["k1"] = make(chan string)

		wg.Add(1)
		go func(c1 <-chan string, c2 chan<- []byte) {
			defer wg.Done()
			message := <-c1 + strconv.Itoa(int(i))
			r2, _ := json.Marshal(message)
			c2 <- r2
			fmt.Printf("Received message %d", i)
		}(c1, c2)
	}

	//fmt.Println(<-m["SS"])

	var messages []string
	j := 0
	for i := range c2 {
		str := i
		var jresp string
		json.Unmarshal([]byte(str), &jresp)
		messages[j] = jresp

		j++
		if j == N {
			fmt.Println("All messages received::", messages)
			wg.Done()
		}
	}

	fmt.Println("Last line of coordinator")
}
