package main

import (
	"strconv"
)

func main() {
	//scanner := bufio.NewScanner(os.Stdin)
	//fmt.Println("Enter the order of the node size on powers of 2:")
	//scanner.Scan()
	ringOrder, _ := strconv.Atoi("32")
	hashID(ringOrder)
}
