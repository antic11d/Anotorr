package main

import (
	"./structs/IO"
	"./structs/Node"
	"encoding/json"
	"fmt"
	"net"

	//"net"
)

var trackerReader = IO.Reader{nil}
var trackerWriter = IO.Writer{nil}

func main() {
	self := Node.InitializeNode()

	//Javljam se trekeru
	tAddr, err := net.ResolveTCPAddr("tcp", "192.168.0.13:9095")
	Node.CheckError(err)
	conn, err := net.DialTCP("tcp",nil, tAddr)
	Node.CheckError(err)

	self.ReqConn = conn

	// Citac i pisac otvoreni ka trekeru za postavjanje requestova
	trackerReader = IO.Reader{self.ReqConn}
	trackerWriter = IO.Writer{self.ReqConn}

	// Javljam sta ja imam od fajlova
	jsonSlice, err := json.Marshal(self.SetMyFiles.ToSlice())
	Node.CheckError(err)

	trackerWriter.Write(string(jsonSlice))

	// Poruka predstavljanja trekera, choose option itd...
	msg := trackerReader.Read()
	fmt.Println(msg)

	var ans string
	_, err = fmt.Scanf("%s", &ans)
	Node.CheckError(err)

	if ans == "D" {
		self.RequestDownload(trackerWriter, trackerReader)
	}

	self.WaitGroup.Add(2)

		go self.ListenTracker()

		go self.ListenPeer()

	self.WaitGroup.Wait()
}