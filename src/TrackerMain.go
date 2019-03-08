package main

import (
	"./structs/File"
	"./structs/Requests"
	"./structs/Tracker"
	"./structs/IO"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"container/list"
)


func handleNode(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from:", conn.RemoteAddr().String())

	var writer = IO.Writer{conn}
	writer.Write("Hello! How are you?\nPlease choose an option(d/u):\n")

	reader := IO.Reader{conn}

	var msg = reader.Read()


	if msg == "d" {
		handleDownload(reader, writer)
	} else if msg == "u" {
		handleUpload(conn)
	} else {
		writer.Write("Choose a valid option")
	}

}

func handleUpload(conn net.Conn) {

	conn.Write([]byte("Give me a info of file you want to upload\n"))

	//u klijenu cemo da statujemo fajl da bismo poslali


	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	checkError(err)

	rootHash := string(recvBuff[:bytesRead])

	tracker.Map[rootHash] = File.File{"Uploaded", 100, 10}
}

func handleDownload(reader IO.Reader, writer IO.Writer) {

	writer.Write("Give me a root hash of file you want and public key\n")

	msg := reader.Read()

	requestFromPeer := Requests.DownloadRequestKey{}

	err := json.Unmarshal([]byte(msg), &requestFromPeer)

	checkError(err)

	fmt.Printf("Got request: %+v\n", requestFromPeer)

	// Hardkodovano maksimalna velicina liste 100
	tracker.DownloadRequests[requestFromPeer] = Requests.DownloadRequest{new(list.List), 0}

	// Ovo ce da ide petljom, prodjem kroz sve u mrezi i svakome se javi da im kazem da neko hoce da skida odredjeni fajl
	// Javljam se svima osim onome ko mi je trazio request!!!!
	conn, err := net.Dial("tcp", "192.168.0.28:9091") // 9091 hardkodovano jer tamo slusa peer
	checkError(err)

	tmpReader := IO.Reader{conn}
	tmpWriter := IO.Writer{conn}

	wrappedObject := Requests.WrappedRequest{requestFromPeer, tracker.DownloadRequests[requestFromPeer]}

	tmpMsg, err := json.Marshal(wrappedObject)
	checkError(err)
	tmpWriter.Write(string(tmpMsg))
	fmt.Printf("Prosao write, i napisao: %+v\n", wrappedObject)

	// Dobijem ip od osobe koja kaze da ima fajl
	peerIP := tmpReader.Read()

	tracker.DownloadRequests[requestFromPeer].CryptedIPs.PushBack(peerIP)


	fmt.Printf("Key: %+v, Value: %+v\n", tracker.DownloadRequests[requestFromPeer], tracker.DownloadRequests[requestFromPeer].CryptedIPs)
}


func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

var tcpAddr, _ = net.ResolveTCPAddr("tcp4", ":9090")
var tracker = Tracker.Tracker{tcpAddr,
							make(map[string]File.File),
				make(map[Requests.DownloadRequestKey]Requests.DownloadRequest),
}

func main() {

	tracker.Map["brena"] = File.File{"Lepa Brena", 100, 10}
	tracker.Map["zorka"] = File.File{"Zorica Brunclik", 100, 10}
	tracker.Map["zvorka"] = File.File{"Zvorinka Milosevic", 100, 10}


	listener, err := net.ListenTCP("tcp", tracker.Addr)

	checkError(err)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error while accepting. Continuing...")
			continue
		}

		go handleNode(conn)
	}
}
