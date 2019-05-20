package Tracker

import (
	"../File"
	"../IO"
	"../Requests"
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set"
	"net"
	"os"
	"sync"
)

type Tracker struct {
	Addr *net.TCPAddr
	Map map[string] *File.File
	DownloadRequests map[Requests.DownloadRequestKey] *Requests.DownloadRequest
	ListOfPeers mapset.Set
	FileList string
	AvailableFiles mapset.Set
	AvailableFileNames mapset.Set
}

var separator = "\n-------------------------------------------------------\n"

func (tracker Tracker) HandleNode(conn *net.TCPConn) {
	defer conn.Close()

	//Inicijalizujemo konekciju ka peer-u
	var writer = IO.Writer{conn}
	var reader = IO.Reader{conn}

	//Ovde da probamo handshake
	callerIP := reader.Read()

	fmt.Println("[HandleNode] handshake: ", callerIP)

	writer.Write("OK")

	tracker.ListOfPeers.Add(callerIP)

	for _, peer := range tracker.ListOfPeers.ToSlice() {
		fmt.Println(peer)
	}

	peersListMsg := reader.Read()

	var sliceList []*File.File
	err := json.Unmarshal([]byte(peersListMsg), &sliceList)
	CheckError(err)

	for _, file := range sliceList {
		tracker.Map[(*file).Name] = file
		tracker.AvailableFiles.Add((*file).Name)
	}

	writer.Write(separator+"Please choose an option D - download (currently supported), S - seeding only:"+separator)

	var option = reader.Read()

	if option == "D" {
		tracker.HandleDownload(callerIP, reader, writer)
	} else if option == "S" {
		writer.Write("Just keep seeding...")
	} else {
		writer.Write("Choose a valid option")
	}
}

func (tracker Tracker) HandleDownload(caller string, reader IO.Reader, writer IO.Writer) {
	// Treba da mu dam spisak svih dostupnih fajlova
	sliceList := tracker.AvailableFiles.ToSlice()

	listOfFiles := ""
	i := 1
	for _, file := range sliceList{
		listOfFiles += fmt.Sprintf("%d. %v\n", i, file)
		i++
	}

	fmt.Println(listOfFiles)

	writer.Write(listOfFiles+"Choose a file from the list:")

	request := reader.Read()
	requestFromPeer := Requests.DownloadRequestKey{}
	err := json.Unmarshal([]byte(request), &requestFromPeer)
	CheckError(err)

	var fInfo *File.File
	fInfo = tracker.Map[requestFromPeer.RootHash]

	fMarshall, err := json.Marshal(fInfo)

	writer.Write(string(fMarshall))

	var helpInt int
	helpInt = 0
	tracker.DownloadRequests[requestFromPeer] = &Requests.DownloadRequest{Requests.Matrix{}, &helpInt}

	// Javljam se svima osim onome ko mi je trazio request!!!!
	var group sync.WaitGroup
	var mutex sync.Mutex
	for i, peer := range tracker.ListOfPeers.ToSlice() {
		peerIP := fmt.Sprintf("%v", peer)
		if peer != caller {
			group.Add(1)

			go tracker.contactPeer(peerIP, i, &requestFromPeer, &group, &mutex)
		} else {
			//fmt.Println("Bato jebiga jedini si peer, " + peerIP)
		}
	}

	group.Wait()

	// ovo da se desi tek kad odblokira wg
	*tracker.DownloadRequests[requestFromPeer].Served = 1

	msgFinal, err := json.Marshal(Requests.WrappedRequest{&requestFromPeer, tracker.DownloadRequests[requestFromPeer]})
	CheckError(err)

	fmt.Printf("%+v \n", tracker.DownloadRequests[requestFromPeer].CryptedIPs)
	writer.Write(string(msgFinal))
}

func (tracker Tracker) contactPeer(pIP string, tID int, requestFromPeer *Requests.DownloadRequestKey, group *sync.WaitGroup, mutex *sync.Mutex)  {
	defer group.Done()
	peerAddr, err := net.ResolveTCPAddr("tcp", pIP+":9096")
	CheckError(err)

	tmpConn, err := net.DialTCP("tcp", nil, peerAddr)
	CheckError(err)

	tmpReader := IO.Reader{tmpConn}
	tmpWriter := IO.Writer{tmpConn}

	wrappedObject := Requests.WrappedRequest{requestFromPeer, tracker.DownloadRequests[*requestFromPeer]}

	tmpMsg, err := json.Marshal(wrappedObject)
	CheckError(err)

	tmpWriter.Write(string(tmpMsg))

	cryptedPIP := make([]byte, 128)

	_, err = tmpReader.Conn.Read(cryptedPIP)
	CheckError(err)

	fmt.Println(cryptedPIP)
	if cryptedPIP[0] != 110 {
		// ovde sinhronizuj tredove
		mutex.Lock()

		tracker.DownloadRequests[*requestFromPeer].CryptedIPs.Arr =
			append(tracker.DownloadRequests[*requestFromPeer].CryptedIPs.Arr, cryptedPIP)

		mutex.Unlock()
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
