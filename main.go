package main

import (
	"fmt"
)

var (
	slack SlackConn
)

func main() {
	fmt.Println("Starting up...")
	slack.Connect()
	telnetListenAndServe()
}
