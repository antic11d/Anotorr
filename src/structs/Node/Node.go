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
	"strings"
	"sync"
	"../File"
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
	ReqConn *net.TCPConn // Konekcija koja se inicijalno ostvaruje za postovanje zahteva
	WaitGroup sync.WaitGroup
	MyFiles map[string] File.File
}

type MsgToNode struct {
	RootHash string
	ChunkNum int64
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
	// Ovde napraviti mapu svih fajlova koje mozes da sidujes
	var wg sync.WaitGroup
	p = &Peer{ID:"idPrvi", PrivateKey:pk, IP: getMyIP(), WaitGroup:wg, MyFiles:make(map[string]File.File)}

	var size int64 = 4391844
	var chunks int64 = 5
	var chunkSize int64 = 1000000

	p.MyFiles["zorka"] = File.File{"Zorica Brunclik - Kada bi me pitali.mp3", &size, &chunks, &chunkSize}
	return p

}

// Hardcoded portovi 9091, 9092

// Cekam da mi se javi treker da mi kaze da neko hoce da skida fajl koji ja potencijalno imam
func (peer Peer) ListenTracker() {
	var tListenAddr, err = net.ResolveTCPAddr("tcp4", ":9091")
	CheckError(err)

	peer.ListenerTracker, err = net.ListenTCP("tcp", tListenAddr)
	CheckError(err)

	for  {
		conn, err := peer.ListenerTracker.AcceptTCP()
		fmt.Println("[ListenTracker] Accepted connection from tracker...")
		if err != nil {
			fmt.Println("Error while accepting connection from tracker, continuing...")
			continue
		}

		go handleTracker(conn)
	}
}

func handleTracker(conn *net.TCPConn) {
	defer conn.Close()
	var tmpReader = IO.Reader{conn}
	var tmpWriter = IO.Writer{conn}

	// Poruka sa objektom koji sadrzi fajl koji neko iz mreze hoce da skida
	msg := tmpReader.Read()

	var wrappedRequest Requests.WrappedRequest
	err := json.Unmarshal([]byte(msg), &wrappedRequest)
	CheckError(err)

	fmt.Printf("[handleTracker] Dobio objekat: %+v\n", wrappedRequest)

	// Ovde treba da prodjem kroz svoji fajl sistem i da vidim da li imam taj fajl, ako imam onda vratim svoj IP trekeru
	tmpWriter.Write(strings.Split(conn.LocalAddr().String(), ":")[0])
}

func (peer Peer) RequestDownload(trackerWriter IO.Writer, trackerReader IO.Reader) {
	// Hardkodovano da hocu download opciju
	trackerWriter.Write("d")

	//ovde ubacujemo dodatni read/write gde meni treker salje listu fajlova i odalte ja sa
	//mojim rootHashom uzimam koliko ima chunkova fajl

	//sada to hardkodujemo

	// STATUS 0 = NIJE SKINUT : STATUS 1 = TRENUTNO SE SKIDA : STATUS 2 = SKINUT
	// Treker trazi root hash i public key, tj DownloadRequestKey
	msg := trackerReader.Read()

	fmt.Println(msg)

	request := Requests.DownloadRequestKey{"zorka", &peer.PrivateKey.PublicKey}
	jsonReq, err := json.Marshal(request)
	CheckError(err)

	// Postujem request za download fajla, DownloadRequestKey
	trackerWriter.Write(string(jsonReq))

	// Sad dobijem informacije o fajlu koji trazim (tek posle ide lista kriptovanih)
	fileInfo := trackerReader.Read()

	var fInfo File.File

	err = json.Unmarshal([]byte(fileInfo), &fInfo)
	CheckError(err)

	fmt.Printf("[RequestDownload] Got file info object: %+v\n", fInfo)

	numOfChunks := *fInfo.Chunks
	var numOfDownloadedChunks int64 = 0

	chunksStatuses := make([]int, numOfChunks)

	var i int64
	for i = 0 ; i < numOfChunks ;i++ {
		chunksStatuses[i] = 0
	}

	// Poruka sa objektom kome sve treba da se javim
	msg = trackerReader.Read()

	completedReq := Requests.WrappedRequest{}
	err = json.Unmarshal([]byte(msg), &completedReq)
	CheckError(err)

	fmt.Printf("[RequestDownload] Treba da se javim svima iz liste: %+v i duzine %+v\n", completedReq.Value.CryptedIPs, len(completedReq.Value.CryptedIPs))

	var downloadWG sync.WaitGroup
	var list = completedReq.Value.CryptedIPs

	f, err := os.Create("/home/antic/Desktop/" + fInfo.Name)

	//Ovde sada cekamo dok se ne skupe svi skinuti cankovi
	tmpSeeder := 0

	var mutex = &sync.Mutex{}

	for i = 0; i < numOfChunks; i++ {
		if numOfDownloadedChunks == numOfChunks {
			break
		}

		if chunksStatuses[i] == 0 {
			// Da se proveri da li se niz salje po referenci tj apdejutuje u funkciji
			chunksStatuses[i] = 1
			downloadWG.Add(1)
			go peer.connectToPeer(list[tmpSeeder], &downloadWG, f, i, chunksStatuses, mutex, &numOfDownloadedChunks)
		}
		tmpSeeder = (tmpSeeder + 1) % len(list)
	}

	CheckError(err)
	downloadWG.Wait()
}

func (peer Peer) ListenPeer() {

	var pListenAddr, err = net.ResolveTCPAddr("tcp4", ":9092")
	CheckError(err)

	peer.ListenerPeer, err = net.ListenTCP("tcp", pListenAddr)
	CheckError(err)

	for  {
		conn, err := peer.ListenerPeer.AcceptTCP()
		fmt.Println("[ListenPeer] Accepted connection from peer...")
		if err != nil {
			fmt.Println("Error while accepting connection from peer, continuing...")
			continue
		}

		go peer.handlePeer(conn)
	}

}

func (peer Peer) handlePeer(conn *net.TCPConn) {
	defer conn.Close()
	var tmpReader= IO.Reader{conn}
	var tmpWriter= IO.Writer{conn}

	msgFromPeer := tmpReader.Read()

	msg := MsgToNode{}
	err := json.Unmarshal([]byte(msgFromPeer), &msg)

	fmt.Printf("[handlePeer] Dobio rootHash: %+v\n", msg.RootHash)

	f, err := os.Open("misc/"+peer.MyFiles[msg.RootHash].Name)
	CheckError(err)
	defer f.Close()

	//tmpWriter.WriteFile("misc/", 0, fInfo.Size())
	// Ovde sajlem File objekat koji downloader prima
	finfo, err := f.Stat()
	CheckError(err)

	tmpWriter.WriteFile(peer.MyFiles[msg.RootHash].Name, msg.ChunkNum, *peer.MyFiles[msg.RootHash].ChunkSize, finfo.Size())
}

func (peer Peer) connectToPeer(IP string, group *sync.WaitGroup, f *os.File, numOfPart int64, chunkStatuses []int, mutex *sync.Mutex, numOfDownloaded *int64) {
	fmt.Printf("[connectToPeer] About to dial: %+v\n", IP)

	rAddr, err := net.ResolveTCPAddr("tcp", IP+":9092")

	conn, err := net.DialTCP("tcp", nil, rAddr)
	CheckError(err)

	tmpReader := IO.Reader{conn}
	tmpWriter := IO.Writer{conn}

	// Hardkodovan root hash fajla koji hocu

	msg := MsgToNode{"zorka", numOfPart}

	msgForSend, err := json.Marshal(msg)
	CheckError(err)

	tmpWriter.Write(string(msgForSend))

	partBytes, size := tmpReader.ReadFile()

	fmt.Printf("wg: %+v\n", group)

	// Ovde mora da se zakljuca fajl pre pisanja
	mutex.Lock()

	_, err = f.WriteAt(partBytes[:size], int64(numOfPart) * size)

	chunkStatuses[numOfPart] = 2

	*numOfDownloaded++

	mutex.Unlock()

	CheckError(err)

	group.Done()
}

