package Node

import (
	"../IO"
	"../Requests"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

type Peer struct {
	ID string
	IP net.IP
	PrivateKey *rsa.PrivateKey
	RootHashes []string
	TrackerListenAddr *net.TCPAddr
	PeerListenAddr *net.TCPAddr
	ListenerTracker *net.TCPListener
	ListenerPeer *net.TCPListener
	ReqConn net.Conn // Konekcija koja se inicijalno ostvaruje za postovanje zahteva
	WaitGroup sync.WaitGroup
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getMyIP() (net.IP) {
	//Ovde ce da se implementira UPNP

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("[getMyIP]IPv4: ", ipv4)
			return ipv4
		}
	}

	return nil
}

func InitializeNode() (p *Peer){
	pk, err := rsa.GenerateKey(rand.Reader, 2048)

	CheckError(err)

	// ID generisati dinamicki!!!
	var wg sync.WaitGroup
	return &Peer{ID:"idPrvi", PrivateKey:pk, IP: getMyIP(), WaitGroup:wg}

}

// Hardcoded portovi 9091, 9092

// Cekam da mi se javi treker da mi kaze da neko hoce da skida fajl koji ja potencijalno imam
func (peer Peer) ListenTracker() {
	var tListenAddr, err = net.ResolveTCPAddr("tcp4", ":9091")
	CheckError(err)

	peer.ListenerTracker, err = net.ListenTCP("tcp", tListenAddr)
	CheckError(err)

	for  {
		conn, err := peer.ListenerTracker.Accept()
		fmt.Println("[ListenTracker] Accepted connection from tracker...")
		if err != nil {
			fmt.Println("Error while accepting connection from tracker, continuing...")
			continue
		}

		go handleTracker(conn)
	}
}

func handleTracker(conn net.Conn) {
	defer conn.Close()
	var tmpReader = IO.Reader{conn}
	var tmpWriter = IO.Writer{conn}

	// Poruka sa objektom koji sadrzi fajl koji neko iz mreze hoce da skida
	msg := tmpReader.Read()

	var wrappedRequest Requests.WrappedRequest

	err := json.Unmarshal([]byte(msg), &wrappedRequest)
	CheckError(err)

	fmt.Printf("[handleTracker] Dobio objekat: %+v\n", wrappedRequest)
	//fmt.Printf("[handleTracker] lista: %+v\n", wrappedRequest.Value.CryptedIPs.Len())

	// Ovde treba da prodjem kroz svoji fajl sistem i da vidim da li imam taj fajl, ako imam onda vratim svoj IP trekeru
	tmpWriter.Write("192.168.0.19")
}

func (peer Peer) RequestDownload(trackerWriter IO.Writer, trackerReader IO.Reader) {
	// Hardkodovano da hocu download opciju
	trackerWriter.Write("d")

	// Treker trazi root hash i public key, tj DownloadRequestKey
	msg := trackerReader.Read()

	fmt.Println(msg)

	request := Requests.DownloadRequestKey{"zorka", &peer.PrivateKey.PublicKey}
	jsonReq, err := json.Marshal(request)

	CheckError(err)

	// Postujem request za download fajla, DownloadRequestKey
	trackerWriter.Write(string(jsonReq))

	// Poruka sa objektom kome sve treba da se javim
	msg = trackerReader.Read()

	completedReq := Requests.WrappedRequest{}
	err = json.Unmarshal([]byte(msg), &completedReq)
	CheckError(err)

	fmt.Printf("[RequestDownload] Treba da se javim svima iz liste: %+v i duzine %+v\n", completedReq.Value.CryptedIPs, len(completedReq.Value.CryptedIPs))

	var downloadWG sync.WaitGroup
	var list = completedReq.Value.CryptedIPs
	downloadWG.Add(len(list))

	// Napravi fajl, i fji posalji fd
	f, err := os.Create("/home/antic/Desktop/proba.txt")
	CheckError(err)
	for i := 0; i < len(list); i++  {
		// ovde ce da se salje i pokazivac na niz koji sadrzi da li je part validan, tu ce biti resen problem ako neko odustane

		go peer.connectToPeer(list[i], &downloadWG, f, 5)
	}

	downloadWG.Wait()
}

func (peer Peer) ListenPeer() {

	var pListenAddr, err = net.ResolveTCPAddr("tcp4", ":9092")
	CheckError(err)

	peer.ListenerPeer, err = net.ListenTCP("tcp", pListenAddr)
	CheckError(err)

	for  {
		conn, err := peer.ListenerPeer.Accept()
		fmt.Println("[ListenPeer] Accepted connection from peer...")
		if err != nil {
			fmt.Println("Error while accepting connection from peer, continuing...")
			continue
		}

		go handlePeer(conn)
	}

}

func handlePeer(conn net.Conn) {
	defer conn.Close()
	var tmpReader= IO.Reader{conn}
	var tmpWriter= IO.Writer{conn}

	rootHash := tmpReader.Read()

	fmt.Printf("[handlePeer] Dobio rootHash: %+v\n", rootHash)

	f, err := os.Open("misc/probaSlika.jpg")
	CheckError(err)

	defer f.Close()

	fInfo, err := f.Stat()

	tmpWriter.WriteFile("misc/probaSlika.jpg", 0, fInfo.Size())
}

func (peer Peer) connectToPeer(IP string, group *sync.WaitGroup, f *os.File, numOfPart int) {
	defer group.Done()

	fmt.Printf("[connectToPeer] About to dial: %+v\n", IP)

	conn, err := net.Dial("tcp", IP + ":9092")
	CheckError(err)

	tmpReader := IO.Reader{conn}
	tmpWriter := IO.Writer{conn}

	// Hardkodovan root hash fajla koji hocu
	tmpWriter.Write("zorka")

	partBytes, partSize := tmpReader.ReadFile()

	// Ovde mora da se zakljuca fajl pre pisanja
	_, err = f.WriteAt(partBytes, int64(partSize*numOfPart))
	CheckError(err)

	fmt.Println("[connectToPeer]", len(partBytes))
}

