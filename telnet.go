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

	var text string

	buf := bufio.NewReader(conn)
	for {
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client " + remoteAddr + " disconnected.")
			break
		}

		fmt.Println(bytes)

		i := 0
		for i < len(bytes) {
			byte := bytes[i]
			if byte == 255 {
				command := bytes[i+1]
				option := bytes[i+2]

				switch command {
				case 254:
					fmt.Printf("IAC DON'T %d\n", option)
					i += 3
				case 253:
					fmt.Printf("IAC DO %d\n", option)
					i += 3
				case 252:
					fmt.Printf("IAC WON'T %d\n", option)
					i += 3
				case 251:
					fmt.Printf("IAC WILL %d\n", option)
					i += 3
				default:
					i++
				}
			} else {
				if byte >= 32 && byte <= 126 {
					text += string(byte)
				}
				i++
			}
		}

		if len(text) > 0 {
			conn.Write([]byte("You said: " + text + "\n"))
			text = ""
		}
	}
}
