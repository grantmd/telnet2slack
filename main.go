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

	// Connect to Slack first because we can't do anything else without it
	slack.Connect()

	// Now that the Slack connection is up, wait for and accept telnet connections
	telnet.ListenAndServe()
}
