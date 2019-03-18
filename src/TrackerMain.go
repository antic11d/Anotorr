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
							make(map[string]File.File),
				make(map[Requests.DownloadRequestKey]*Requests.DownloadRequest),
				make([]string, 0),
}

func main() {
	tracker.Map["brena"] = File.File{"Lepa Brena", 100, 10}
	tracker.Map["zorka"] = File.File{"Zorica Brunclik", 100, 10}
	tracker.Map["zvorka"] = File.File{"Zvorinka Milosevic", 100, 10}


	listener, err := net.ListenTCP("tcp", tracker.Addr)

	Tracker.CheckError(err)

	for {
		conn, err := listener.AcceptTCP()
		fmt.Println("[TrackerMain] Got call from", strings.Split(conn.RemoteAddr().String(), ":")[0])

		//tracker.ListOfPeers = append(tracker.ListOfPeers, strings.Split(conn.RemoteAddr().String(), ":")[0])

		tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.162.98")
		//tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.151.148")

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go tracker.HandleNode(conn)
	}
}
