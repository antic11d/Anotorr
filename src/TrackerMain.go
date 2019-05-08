	package main

import (
	"./structs/File"
	"./structs/Requests"
	"./structs/Tracker"
	"fmt"
	"gitlab.com/NebulousLabs/go-upnp"
	"net"
	"strings"
)

var separator = "\n--------------------------------------------\n"
var tcpAddr, _ = net.ResolveTCPAddr("tcp4", ":9095")
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

	d, err := upnp.Discover()
	Tracker.CheckError(err)

	// discover external IP
	ip, err := d.ExternalIP()
	Tracker.CheckError(err)
	fmt.Println(separator+"Your external IP is:" + ip + separator)

	err = d.Forward(9095, "upnp goTorr 2")
	Tracker.CheckError(err)

	listener, err := net.ListenTCP("tcp", tracker.Addr)

	Tracker.CheckError(err)

	for {
		conn, err := listener.AcceptTCP()

		caller := strings.Split(conn.RemoteAddr().String(), ":")[0]
		fmt.Println("[TrackerMain] Got call from", caller)

		tracker.ListOfPeers = append(tracker.ListOfPeers, caller)
		fmt.Println("Added peer " + caller)

		//tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.162.98")
		//tracker.ListOfPeers = append(tracker.ListOfPeers, "192.168.1.106")
		//tracker.ListOfPeers = append(tracker.ListOfPeers, "10.0.155.169")

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go tracker.HandleNode(conn)
	}
}
