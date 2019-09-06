package main

import (
	"fmt"
)

var (
	slack  SlackConn
	telnet TelnetServer
)

func main() {
	fmt.Println("Starting up...")
	slack.Connect()
	telnet.ListenAndServe()
}
