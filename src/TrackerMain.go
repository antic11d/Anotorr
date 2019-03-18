	package main

import (
	"./structs/File"
	"./structs/Requests"
	"./structs/Tracker"
	"fmt"
	"net"
	"strings"
)

var tcpAddr, _ = net.ResolveTCPAddr("tcp4", ":9090")
var tracker = Tracker.Tracker{tcpAddr,
							make(map[string] *File.File),
				make(map[Requests.DownloadRequestKey]*Requests.DownloadRequest),
				make([]string, 0),
}

func main() {
	//tracker.Map["brena"] = &File.File{"Lepa Brena", 100, 10, 123}
	var size int64 = 4391844
	var chunks int64 = 5
	var chunkSize int64 = 1000000
	tracker.Map["zorka"] = &File.File{"zorka.mp3", &size, &chunks, &chunkSize}
	//tracker.Map["zvorka"] = &File.File{"Zvorinka Milosevic", 100, 10, 123}

	listener, err := net.ListenTCP("tcp", tracker.Addr)

	Tracker.CheckError(err)

	for {
		conn, err := listener.AcceptTCP()
		fmt.Println("[TrackerMain] Got call from", strings.Split(conn.RemoteAddr().String(), ":")[0])

		//tracker.ListOfPeers = append(tracker.ListOfPeers, strings.Split(conn.RemoteAddr().String(), ":")[0])

		//tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.162.98")
		tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.162.98")
		tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.155.169")

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go tracker.HandleNode(conn)
	}
}
