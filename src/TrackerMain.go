	package main

	import (
		"./structs/File"
		"./structs/Requests"
		"./structs/Tracker"
		"fmt"
		"github.com/deckarep/golang-set"
		"net"
		"strings"
	)

var separator = "\n--------------------------------------------\n"
var tcpAddr, _ = net.ResolveTCPAddr("tcp4", ":9095")
var tracker = Tracker.Tracker{tcpAddr,
							make(map[string] *File.File),
				make(map[Requests.DownloadRequestKey]*Requests.DownloadRequest),
				mapset.NewSet(),
				"",
				mapset.NewSet(),
				mapset.NewSet(),
}

func main() {/*
	d, err := upnp.Discover()
	Tracker.CheckError(err)

	// discover external IP
	ip, err := d.ExternalIP()
	Tracker.CheckError(err)
	fmt.Println(separator+"Your external IP is:" + ip + separator)

	err = d.Forward(9095, "upnp goTorr 2")
	Tracker.CheckError(err)*/

	listener, err := net.ListenTCP("tcp", tracker.Addr)

	Tracker.CheckError(err)

	for {
		conn, err := listener.AcceptTCP()

		caller := strings.Split(conn.RemoteAddr().String(), ":")[0]
		fmt.Println("[TrackerMain] Got call from", caller, " [cvoriste...]")

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go tracker.HandleNode(conn)
	}
}
