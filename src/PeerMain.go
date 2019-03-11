package main

import (
	"./structs/IO"
	"./structs/Node"
	"./structs/Requests"
	"encoding/json"
	"fmt"
	"net"
)

var trackerReader = IO.Reader{nil}
var trackerWriter = IO.Writer{nil}

func main() {
	var selfNode = Node.InitializeNode()

	fmt.Printf("Hello, my name is: %+v\n", selfNode)

	conn, err := net.Dial("tcp", "127.0.0.1:9090")

	Node.CheckError(err)

	trackerReader = IO.Reader{conn}
	trackerWriter = IO.Writer{conn}

	msg := trackerReader.Read()

	fmt.Println(msg)

	downloadFile(conn, selfNode)

	fmt.Println("About to listen on port for tracker info...")

	selfNode.ListenTracker()
}

func downloadFile(conn net.Conn, peer *Node.Peer) {
	//fmt.Scanln()
	trackerWriter.Write("d")

	msg := trackerReader.Read()

	fmt.Println(msg)

	request := Requests.DownloadRequestKey{"zorka", &peer.PrivateKey.PublicKey}
	jsonReq, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	trackerWriter.Write(string(jsonReq))
}


