package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	telnetPort = 23
)

// TelnetServer listens for incoming telnet connections
type TelnetServer struct {
	ln net.Listener
}

// ListenAndServe waits for incoming telnet connections and then dispatches them
func (t *TelnetServer) ListenAndServe() (err error) {
	t.ln, err = net.Listen("tcp", fmt.Sprintf(":%d", telnetPort))
	if err != nil {
		return err
	}

	fmt.Printf("Listening on port %d\n", telnetPort)

	for {
		conn, err := t.ln.Accept()
		if err != nil {
			return err
		}

		go handleTelnetConnection(conn)
	}

	return nil
}

func handleTelnetConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("New connection from " + remoteAddr)

	writeTelnetOutput(conn, "Welcome to telnet2slack. What is your name?\r\n")

	var name string
	var input string

	buf := bufio.NewReader(conn)
	for {
		bytes, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println("Client " + remoteAddr + " disconnected.")
			break
		}

		input = readTelnetInput(bytes)

		if len(input) > 0 {
			if name == "" {
				name = input
				writeTelnetOutput(conn, "Hello "+input+"!\r\n")
				slack.SendMessage(fmt.Sprintf("%s is here!", name))
			} else {
				writeTelnetOutput(conn, name+": "+input+"\r\n")
				slack.SendMessage(fmt.Sprintf("%s: %s", name, input))
			}
		}
	}
}

func readTelnetInput(bytes []byte) string {
	var input string
	fmt.Println(bytes)

	i := 0
	for i < len(bytes) {
		byte := bytes[i]
		if byte == 255 {
			command := bytes[i+1]
			option := bytes[i+2]

			switch command {
			// https://www.iana.org/assignments/telnet-options/telnet-options.xhtml
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
				input += string(byte)
			}
			i++
		}
	}

	return input
}

func writeTelnetOutput(conn net.Conn, output string) {
	conn.Write([]byte(output))
}
