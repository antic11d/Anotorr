package Tracker

import (
	"../File"
	"../Requests"
	"../IO"
	//"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Tracker struct {
	Addr *net.TCPAddr
	Map map[string]File.File
	DownloadRequests map[Requests.DownloadRequestKey] *Requests.DownloadRequest
}


func (tracker Tracker) HandleNode(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from:", conn.RemoteAddr().String())

	//Inicijalizujemo konekciju ka peer-u
	var writer = IO.Writer{conn}
	var reader = IO.Reader{conn}

	writer.Write("Hello! How are you?\nPlease choose an option(d/u):\n")

	var option = reader.Read()

	if option == "d" {
		tracker.HandleDownload(reader, writer)
	} else if option == "u" {
		tracker.HandleUpload(conn)
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

	// Hardkodovano maksimalna velicina liste 100
	var helpInt int
	helpInt = 0
	tracker.DownloadRequests[requestFromPeer] = &Requests.DownloadRequest{make([]string, 0), &helpInt}

	// Ovo ce da ide petljom, prodjem kroz sve u mrezi i svakome se javi da im kazem da neko hoce da skida odredjeni fajl
	// Javljam se svima osim onome ko mi je trazio request!!!!
	tmpConn, err := net.Dial("tcp", "10.0.162.98:9091") // 9091 hardkodovano jer tamo slusa peer
	CheckError(err)

	tmpReader := IO.Reader{tmpConn}
	tmpWriter := IO.Writer{tmpConn}

	wrappedObject := Requests.WrappedRequest{&requestFromPeer, tracker.DownloadRequests[requestFromPeer]}

	tmpMsg, err := json.Marshal(wrappedObject)
	CheckError(err)

	fmt.Printf("[HandleDownload] Poslao %+v, objekat: %+v\n", tmpWriter.Conn.RemoteAddr(), wrappedObject)
	tmpWriter.Write(string(tmpMsg))

	// Dobijem ip od osobe koja kaze da ima fajl
	// Ovde ce IP biti kodiran normalno
	peerIP := tmpReader.Read()
	fmt.Println("[HandleDownload] dobio ip:" + peerIP + " od peera: " + tmpReader.Conn.RemoteAddr().String())

	// Dodajemo ga u listu koju treba poslati onome ko je trazio request
	// Ovde ce morati i sinhronizacija, da se zakljuca mapa
	fmt.Printf("[HandleDownload] Dodajem kriptovani IP u listu...\n")
	//tracker.DownloadRequests[requestFromPeer].CryptedIPs.PushBack(peerIP)
	tracker.DownloadRequests[requestFromPeer].CryptedIPs =
					append(tracker.DownloadRequests[requestFromPeer].CryptedIPs, peerIP)

	// Nakon sto se "napuni" lista kriptovanih, posalji onome ko je trazio ceo objekat
	fmt.Printf("[HandleDownload] Objekat koji treba poslati peer-u da moze da se javi kome treba itd...\nKey: %+v, duzina liste: %+v\n",
		tracker.DownloadRequests[requestFromPeer], len(tracker.DownloadRequests[requestFromPeer].CryptedIPs))

	// sto ne moze???
	//served *int := 1
	*tracker.DownloadRequests[requestFromPeer].Served = 1

	fmt.Println("Ovo je served kod nas\n", tracker.DownloadRequests[requestFromPeer].Served)

	msgFinal, err := json.Marshal(Requests.WrappedRequest{&requestFromPeer, tracker.DownloadRequests[requestFromPeer]})
	CheckError(err)

	// Lista prazna???
	fmt.Println("[HandleDownload] msgFinal:" + string(msgFinal))
	fmt.Printf("%+v \n", tracker.DownloadRequests[requestFromPeer].CryptedIPs)
	writer.Write(string(msgFinal))
}

func (tracker Tracker) HandleUpload(conn net.Conn) {

	conn.Write([]byte("Give me a info of file you want to upload\n"))

	//u klijenu cemo da statujemo fajl da bismo poslali


	recvBuff := make([]byte, 2048)

	bytesRead, err := conn.Read(recvBuff)

	CheckError(err)

	rootHash := string(recvBuff[:bytesRead])

	tracker.Map[rootHash] = File.File{"Uploaded", 100, 10}
}


func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
