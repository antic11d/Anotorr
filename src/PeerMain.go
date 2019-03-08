package main

import (
	"./structs/File"
	"./structs/Node"
	"./structs/IO"
	"encoding/json"
	"fmt"
	"net"
)

var trackerReader = IO.Reader{nil}
var trackerWriter = IO.Writer{nil}


func main() {

	var selfNode = Node.InitializeNode()

	fmt.Println(selfNode)

	conn, err := net.Dial("tcp", "127.0.0.1:9090")

	Node.CheckError(err)

	trackerReader = IO.Reader{conn}
	trackerWriter = IO.Writer{conn}

	msg := trackerReader.Read()

	fmt.Println(msg)

	downloadFile(conn)





}

func downloadFile(conn net.Conn) {

	trackerWriter.Write("d")

	msg := trackerReader.Read()

	fmt.Println(msg)

	trackerWriter.Write("zorka")

	msg = trackerReader.Read()

	var f File.File

	err := json.Unmarshal([]byte(msg), &f)
	Node.CheckError(err)

	fmt.Printf("%+v\n", f)
}


