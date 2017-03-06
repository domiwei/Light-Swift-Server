package main

import "fmt"
import "light-swift-server/swifttest"

func main() {
	_, err := swifttest.NewSwiftServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server Start!")
	select {}
}
