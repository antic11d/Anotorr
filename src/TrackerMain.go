package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"./structs/Tracker"
	"./structs/File"
)


func handleNode(conn net.Conn) {
	defer conn.Close()

	_, err := conn.Write([]byte("Hello! How are you?\nPlease choose an option(d/u):\n"))

	checkError(err)

	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	checkError(err)

	str := string(recvBuff[:bytesRead-2])

	if str == "d" {
		handleDownload(conn)
	} else if str == "u" {
		handleUpload(conn)
	} else {
		conn.Write([]byte("Choose a valid option"))
	}

}

func handleUpload(conn net.Conn) {

	conn.Write([]byte("Good job"))
}

func handleDownload(conn net.Conn) {

	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	checkError(err)

	str := string(recvBuff[:bytesRead-2])

	for k, v := range tracker.Map {

		if k == str {

			msg, err := json.Marshal(v)

			checkError(err)

			conn.Write([]byte(msg))
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
var tracker = Tracker.Tracker{tcpAddr, make(map[string]File.File)}

func main() {

	tracker.Map["brena"] = File.File{"Lepa Brena", 100, 10}
	tracker.Map["zorka"] = File.File{"Zorica Brunclik", 100, 10}
	tracker.Map["zvorka"] = File.File{"Zvorinka Milosevic", 100, 10}


	listener, err := net.ListenTCP("tcp", tracker.Addr)

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
