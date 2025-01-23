package main

import "fmt"

func main() {
	symbol := []byte{0xE2, 0x96, 0xBA} // UTF-8 код для символа "►"
	fmt.Println(string(symbol))
}
