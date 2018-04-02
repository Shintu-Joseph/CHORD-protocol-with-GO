package main

import (
	"fmt"
	"sync"
)

func main() {

	// fmt.Println("Enter the number of nodes to spawn:: ")
	// fmt.Scan(&N)
	// if N == 0 || N < 0 {
	// 	fmt.Println("Please enter positive number of nodes")
	// 	os.Exit(0)
	// }

	fmt.Println("Running Coordinator")

	var wg sync.WaitGroup
	wg.Add(1)

	go coordinator(generateRandomID(20), &wg)

	//scanner := bufio.NewScanner(os.Stdin)
	//fmt.Println("Enter the order of the node size on powers of 2:")
	//scanner.Scan()
	wg.Wait()

}
