package Tracker

import (
	"../File"
	"../Requests"
	"../IO"
	"sync"

	//"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Tracker struct {
	Addr *net.TCPAddr
	Map map[string] *File.File
	DownloadRequests map[Requests.DownloadRequestKey] *Requests.DownloadRequest
	ListOfPeers []string
}

var separator = "\n================================================\n"

func (tracker Tracker) HandleNode(conn *net.TCPConn) {
	defer conn.Close()

	fmt.Println("Accepted connection from:", conn.RemoteAddr().String())

	//Inicijalizujemo konekciju ka peer-u
	var writer = IO.Writer{conn}
	var reader = IO.Reader{conn}

	writer.Write(separator+"Hello! How are you?\nPlease choose an option D - download (currently supported):"+separator)

	var option = reader.Read()

	if option == "D" {
		tracker.HandleDownload(reader, writer)
	} else if option == "u" {
		//tracker.HandleUpload(conn)
	} else {
		writer.Write("Choose a valid option")
	}

}

func (tracker Tracker) HandleDownload(reader IO.Reader, writer IO.Writer) {
	writer.Write("Give me a root hash of file you want and public key\n")

	request := reader.Read()

	requestFromPeer := Requests.DownloadRequestKey{}

	err := json.Unmarshal([]byte(request), &requestFromPeer)
	CheckError(err)

	fmt.Printf("Got request: %+v from %+v\n", requestFromPeer, reader.Conn.RemoteAddr())

	var fInfo *File.File

	fInfo = tracker.Map[requestFromPeer.RootHash]

	fMarshall, err := json.Marshal(fInfo)

	writer.Write(string(fMarshall))

	// Hardkodovano maksimalna velicina liste 100
	var helpInt int
	helpInt = 0
	tracker.DownloadRequests[requestFromPeer] = &Requests.DownloadRequest{make([]string, 0), &helpInt}

	// Javljam se svima osim onome ko mi je trazio request!!!!
	var group sync.WaitGroup
	var mutex sync.Mutex
	for i, peer := range tracker.ListOfPeers {
		group.Add(1)
		go tracker.contactPeer(peer, i, &requestFromPeer, &group, &mutex)
	}

	group.Wait()

	// ovo da se desi tek kad odblokira wg
	*tracker.DownloadRequests[requestFromPeer].Served = 1

	msgFinal, err := json.Marshal(Requests.WrappedRequest{&requestFromPeer, tracker.DownloadRequests[requestFromPeer]})
	CheckError(err)

	fmt.Println("[HandleDownload] msgFinal:" + string(msgFinal))
	fmt.Printf("%+v \n", tracker.DownloadRequests[requestFromPeer].CryptedIPs)
	writer.Write(string(msgFinal))
}

func (tracker Tracker) contactPeer(pIP string, tID int, requestFromPeer *Requests.DownloadRequestKey, group *sync.WaitGroup, mutex *sync.Mutex)  {
	defer group.Done()
	peerAddr, err := net.ResolveTCPAddr("tcp", pIP+":9091")
	CheckError(err)

	tmpConn, err := net.DialTCP("tcp", nil, peerAddr) // 9091 hardkodovano jer tamo slusa peer
	CheckError(err)

	tmpReader := IO.Reader{tmpConn}
	tmpWriter := IO.Writer{tmpConn}

	wrappedObject := Requests.WrappedRequest{requestFromPeer, tracker.DownloadRequests[*requestFromPeer]}

	tmpMsg, err := json.Marshal(wrappedObject)
	CheckError(err)

	fmt.Printf("[HandleDownload] %d-tom iz liste Poslao %+v, objekat: %+v\n", tID, tmpWriter.Conn.RemoteAddr(), wrappedObject)
	tmpWriter.Write(string(tmpMsg))

	peerIP := tmpReader.Read()
	fmt.Println("[HandleDownload] dobio ip:" + peerIP + " od peera: " + tmpReader.Conn.RemoteAddr().String())

	fmt.Printf("[HandleDownload] Dodajem kriptovani IP u listu koju cu da posaljem kad se napuni...\n")

	// ovde sinhronizuj tredove
	mutex.Lock()
	tracker.DownloadRequests[*requestFromPeer].CryptedIPs =
		append(tracker.DownloadRequests[*requestFromPeer].CryptedIPs, peerIP)
	mutex.Unlock()
}
/*
func (tracker Tracker) HandleUpload(conn net.Conn) {

	conn.Write([]byte("Give me a info of file you want to upload\n"))

	//u klijenu cemo da statujemo fajl da bismo poslali


	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	CheckError(err)

	//rootHash := string(recvBuff[:bytesRead])

	//tracker.Map[rootHash] = File.File{"Uploaded", 100, 10}
}*/


func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
