package main

import (
	"bufio"
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

	fmt.Printf("Listening on port %d\n", telnetPort)

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
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("New connection from " + remoteAddr)

	buf := bufio.NewReader(conn)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Client " + remoteAddr + " disconnected.")
			break
		}

		conn.Write([]byte("You said: " + line))
	}
}
