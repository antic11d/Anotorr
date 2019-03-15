package main

import (
	"./structs/IO"
	"./structs/Node"
	"fmt"
	"net"
)

var trackerReader = IO.Reader{nil}
var trackerWriter = IO.Writer{nil}

func main() {
	var self = Node.InitializeNode()

	fmt.Printf("[PeerMain] Hello, my name is: %+v\n", self)

	//Javljam se trekeru. Hardkodovan localhost
	tAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9090")
	conn, err := net.DialTCP("tcp",nil, tAddr)
	Node.CheckError(err)

	self.ReqConn = conn

	// Citac i pisac otvoreni ka trekeru za postavjanje requestova
	trackerReader = IO.Reader{self.ReqConn}
	trackerWriter = IO.Writer{self.ReqConn}

	// Poruka predstavljanja trekera, choose option itd...
	msg := trackerReader.Read()
	fmt.Println(msg)

	self.RequestDownload(trackerWriter, trackerReader)

	fmt.Println("About to listen on port for tracker info...")

	self.WaitGroup.Add(2)


	go self.ListenTracker()

	go self.ListenPeer()

	self.WaitGroup.Wait()
}