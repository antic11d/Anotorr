package main

import (
	"./structs/File"
	"./structs/Requests"
	"./structs/Tracker"
	"./structs/IO"
	"encoding/json"
	"fmt"
	"net"
	"os"
)


func handleNode(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from:", conn.RemoteAddr().String())

	var writer IO.Writer = IO.Writer{conn}
	writer.Write("Hello! How are you?\nPlease choose an option(d/u):\n")

	reader := IO.Reader{conn}

	var msg = reader.Read()


	if msg == "d" {
		handleDownload(conn)
	} else if msg == "u" {
		handleUpload(conn)
	} else {
		writer.Write("Choose a valid option")
	}

}

func handleUpload(conn net.Conn) {

	conn.Write([]byte("Give me a info of file you want to upload\n"))

	//u klijenu cemo da statujemo fajl da bismo poslali


	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	checkError(err)

	rootHash := string(recvBuff[:bytesRead])

	tracker.Map[rootHash] = File.File{"Uploaded", 100, 10}

}

func handleDownload(conn net.Conn) {

	conn.Write([]byte("Give me a root hash of file you want\n"))

	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	checkError(err)

	str := string(recvBuff[:bytesRead])

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
var tracker = Tracker.Tracker{tcpAddr,
							make(map[string]File.File),
				make(map[Requests.DownloadRequestKey]Requests.DownloadRequest),
}

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
