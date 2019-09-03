package main

import (
	"fmt"
	"net"
)

var (
	telnetPort = 23
)

func telnetListenAndServe() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", telnetPort))
	if err != nil {
		fmt.Println("Listen error:")
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:")
			fmt.Println(err)
		}

		go handleTelnetConnection(conn)
	}
}

func handleTelnetConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection from: " + conn.RemoteAddr().String())
}
