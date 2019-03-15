	package main

import (
	"./structs/File"
	"./structs/Requests"
	"./structs/Tracker"
	"fmt"
	"net"
)

var tcpAddr, _ = net.ResolveTCPAddr("tcp4", ":9090")
var tracker = Tracker.Tracker{tcpAddr,
							make(map[string]File.File),
				make(map[Requests.DownloadRequestKey]*Requests.DownloadRequest),
}

func main() {
	tracker.Map["brena"] = File.File{"Lepa Brena", 100, 10}
	tracker.Map["zorka"] = File.File{"Zorica Brunclik", 100, 10}
	tracker.Map["zvorka"] = File.File{"Zvorinka Milosevic", 100, 10}


	listener, err := net.ListenTCP("tcp", tracker.Addr)

	Tracker.CheckError(err)

	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go tracker.HandleNode(conn)
	}
}
