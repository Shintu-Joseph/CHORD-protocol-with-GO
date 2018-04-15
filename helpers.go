package main

func convertArrayToByteArray(array []HashKey) []byte {
	byteArray := make([][]byte, len(array))
	for i, v := range array {
		byteArray[i] = []byte(v)
	}
	return byteArray
}
