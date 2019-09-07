package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	slack  SlackConn
	telnet TelnetServer
)

func main() {
	// Setup signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Quit requested")
		telnet.Close()
		slack.Close()
		os.Exit(1)
	}()

	fmt.Println("Starting up...")

	// Connect to Slack first because we can't do anything else without it
	slack.Connect()

	// Now that the Slack connection is up, wait for and accept telnet connections
	telnet.ListenAndServe()
}
