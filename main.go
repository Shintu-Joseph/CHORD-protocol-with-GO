package main

import (
	"fmt"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func main() {

	var N int
	// fmt.Println("Enter the number of nodes to spawn:: ")
	// fmt.Scan(&N)
	// if N == 0 || N < 0 {
	// 	fmt.Println("Please enter positive number of nodes")
	// 	os.Exit(0)
	// }
	N = 3

	fmt.Println("Running Coordinator")
	wg.Add(1)
	go coordinator(N)
	wg.Wait()
	fmt.Println("Exited Coordinator")
	//scanner := bufio.NewScanner(os.Stdin)
	//fmt.Println("Enter the order of the node size on powers of 2:")
	//scanner.Scan()
	ringOrder, _ := strconv.Atoi("32")
	hashID(ringOrder)
}
