package main

import (
	"fmt"
	"net"
	"os"
)

type Tracker struct {
	addr *net.TCPAddr
	listOfItems []string
}

func handleNode(conn net.Conn) {
	defer conn.Close()

	_, err := conn.Write([]byte("Hello! How are you?\nPlease choose a file:\n"))

	checkError(err)

	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	checkError(err)

	str := string(recvBuff[:bytesRead-2])

	for _, itemName := range tracker.listOfItems {
		if itemName == str {
			_, err := conn.Write([]byte(itemName))
			checkError(err)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

var tcpAddr, _ = net.ResolveTCPAddr("tcp4", ":9090")
var tracker = Tracker{tcpAddr, []string{"Dimitrije", "Stefan", "Andrija"}}

func main() {
	listener, err := net.ListenTCP("tcp", tracker.addr)

	checkError(err)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go handleNode(conn)
	}
}
