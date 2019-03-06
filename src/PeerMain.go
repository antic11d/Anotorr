package main

import (
	"./structs/File"
	"./structs/Node"
	"encoding/json"
	"fmt"
	"net"
)

func main() {

	var selfNode = Node.InitializeNode()

	fmt.Println(selfNode)

	conn, err := net.Dial("tcp", "127.0.0.1:9090")

	Node.CheckError(err)

	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	Node.CheckError(err)

	trackerWelcomeMsg := string(recvBuff[:bytesRead-1])

	fmt.Println(trackerWelcomeMsg)


	// Ovde da dodamo unos za standardni ulaz kad dodje upload u opciju

	//reader := bufio.NewReader(os.Stdin)

	//text, _ := reader.ReadString('\n')

	//if text == "d" {
		downloadFile(conn)
	//}




}

func downloadFile(conn net.Conn) {

	conn.Write([]byte("d"))

	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	Node.CheckError(err)

	trackerMsg := string(recvBuff[:bytesRead-1])

	fmt.Println(trackerMsg)


	//reader := bufio.NewReader(os.Stdin)

	//rootHash, _ := reader.ReadString('\n')

	conn.Write([]byte("zorka"))

	bytesRead, err = conn.Read(recvBuff)

	Node.CheckError(err)

	readFile := recvBuff[:bytesRead]

	var f File.File

	err = json.Unmarshal(readFile, &f)
	Node.CheckError(err)

	fmt.Println(f)
}


