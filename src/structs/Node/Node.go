package Node

import (
	"../File"
	"../IO"
	"../MerkleTree"
	"../Requests"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set"
	"gitlab.com/NebulousLabs/go-upnp"
	"net"
	"os"
	"path/filepath"
	"sync"
)

type Peer struct {
	ID string
	IP string
	PrivateKey *rsa.PrivateKey
	RootHashes []string
	TrackerListenAddr *net.TCPAddr
	PeerListenAddr *net.TCPAddr
	ListenerTracker *net.TCPListener
	ListenerPeer *net.TCPListener
	ReqConn *net.TCPConn // Konekcija koja se inicijalno ostvaruje za postovanje zahteva
	WaitGroup sync.WaitGroup
	MyFolderPath string
	MyFiles map[string] File.File
	MyTrees map[string] MerkleTree.Merkle //Za svaki root hash ja cuvam merkle stablo za njega
	SetMyfNames mapset.Set
	SetMyFiles mapset.Set
	LocalAddr string
}

type MsgToNode struct {
	RootHash string
	ChunkNum int64
}

var separator = "\n-------------------------------------------------------\n"
var FOLDER_PATH = ""

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getMyIP() (string) {
	fmt.Println(separator+"Getting external ip and mapping ports.\nPlease wait for a couple of seconds..."+separator)

	d, err := upnp.Discover()
	CheckError(err)

	// Hvatanje externe ip
	ip, err := d.ExternalIP()
	CheckError(err)
	fmt.Println(separator+"Your external IP is:" + ip + separator)

	// port forwarding
	err = d.Forward(9093, "upnp goTorr 1")
	CheckError(err)

	err = d.Forward(9091, "upnp goTorr 2")
	CheckError(err)

	return ip
}

func getLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("[getLocalIP]IPv4: ", ipv4)
			return ipv4.String()
		}
	}

	return ""
}

func InitializeNode() (p *Peer){
	//Inicijalizaciju Merkle stabla isto uradi ovde
	pk, err := rsa.GenerateKey(rand.Reader, 2048)

	CheckError(err)

	fmt.Println(separator+"Hello node :)\nWhat is your name? RUN AS SUDO!!!"+separator)
	var name string
	//_, err = fmt.Scanf("%s", &name)
	CheckError(err)
	name = "Sta_god"

	var wg sync.WaitGroup
	p = &Peer{ID:name, PrivateKey:pk, IP: getMyIP(), WaitGroup:wg, LocalAddr: getLocalIP()}
	p.MyFiles, p.SetMyfNames, p.SetMyFiles = initListOfFiles()

	p.MyFolderPath = FOLDER_PATH

	/*
	for k, v := range p.MyFiles {
		fmt.Printf("%+v -> ", k)
		fmt.Printf("%+v %+v %+v\n", *v.Size, *v.ChunkSize, *v.Chunks)
	}*/

	return p
}

// moja putanja: /home/goTorr_files
func checkFolder() string {
	/*
	// Malo hardkoda... Trazi se od korisnika da unese putanju do foldera sa fajlovima, inace ima dosta probelma
	// sa pravima pristupa, hijerarhijom unutar /home foldera itd...

	fmt.Println(separator+"Give me a path to goTorr_files folder: (format: /absolute/path/to/folder/goTorr_files)")
	fmt.Println("In case you haven't made it yet type N, mkdir and then start app again, thank you!"+separator)

	var path string
	_, err := fmt.Scanf("%s", &path)
	CheckError(err)

	if path == "N" {
		os.Exit(1)
	}

	finfo, err := os.Stat(path)
	CheckError(err)

	if finfo.IsDir() && finfo.Name() == "goTorr_files" {
		fmt.Println(separator+"All good! Welcome to goTorr community :)"+separator)
	}

	FOLDER_PATH = path

	//return path*/
	return "/home/goTorr_files"
}

func initListOfFiles() (map[string] File.File, mapset.Set, mapset.Set) {
	files := make(map[string]File.File)
	set := mapset.NewSet()
	fSet := mapset.NewSet()
	path := checkFolder()
	FOLDER_PATH = path

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			var fSize = info.Size()

			// Ovde dodati dinamicko deljenje u odnosu na velicinu fajla
			var chunks int64 = 5

			var chunkSize = fSize / chunks

			currFile := File.File{info.Name(), &fSize, &chunks, &chunkSize}
			files[info.Name()] = currFile
			set.Add(info.Name())
			fSet.Add(&currFile)
		}

		return nil
	})
	CheckError(err)

	return files, set, fSet
}

//portovi 9091 = 50335, 9092 = 50336

// Cekam da mi se javi treker da mi kaze da neko hoce da skida fajl koji ja potencijalno imam
func (peer Peer) ListenTracker() {
	var tListenAddr, err = net.ResolveTCPAddr("tcp4", peer.LocalAddr+":9091")
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

		go peer.handleTracker(conn)
	}
}

func (peer Peer) handleTracker(conn *net.TCPConn) {
	defer conn.Close()
	var tmpReader = IO.Reader{conn}
	var tmpWriter = IO.Writer{conn}

	// Poruka sa objektom koji sadrzi fajl koji neko iz mreze hoce da skida
	msg := tmpReader.Read()

	var wrappedRequest Requests.WrappedRequest
	err := json.Unmarshal([]byte(msg), &wrappedRequest)
	CheckError(err)

	fmt.Printf("[handleTracker] Dobio objekat: %+v\n", wrappedRequest)

	// Ako imam fajl javljam svoj IP da bi mi se downloader javio
	if peer.SetMyfNames.Contains(wrappedRequest.Key.RootHash) {
		fmt.Println("[handletracker] dajem svoj ip: "+peer.IP)
		tmpWriter.Write(peer.IP)
	}
}

func (peer Peer) RequestDownload(trackerWriter IO.Writer, trackerReader IO.Reader) {
	trackerWriter.Write("D")

	filesList := trackerReader.Read()

	fmt.Println(separator+"Avaliable files:\n"+filesList+separator)
	// STATUS 0 = NIJE SKINUT, STATUS 1 = TRENUTNO SE SKIDA, STATUS 2 = SKINUT

	var fname string
	_, err := fmt.Scanf("%s\n", &fname)
	CheckError(err)

	request := Requests.DownloadRequestKey{fname, &peer.PrivateKey.PublicKey}
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
	for i = 0; i < numOfChunks; i++ {
		chunksStatuses[i] = 0
	}

	// Poruka sa objektom kome sve treba da se javim
	msg := trackerReader.Read()

	completedReq := Requests.WrappedRequest{}
	err = json.Unmarshal([]byte(msg), &completedReq)
	CheckError(err)

	fmt.Printf("[RequestDownload] Treba da se javim svima iz liste: %+v i duzine %+v\n", completedReq.Value.CryptedIPs, len(completedReq.Value.CryptedIPs))

	var downloadWG sync.WaitGroup
	var list = completedReq.Value.CryptedIPs

	f, err := os.Create(peer.MyFolderPath + "/" + fInfo.Name)

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
			go peer.connectToPeer(fname, list[tmpSeeder], &downloadWG, f, i, chunksStatuses, mutex, &numOfDownloadedChunks)
		}
		tmpSeeder = (tmpSeeder + 1) % len(list)
	}

	CheckError(err)
	downloadWG.Wait()
}

func (peer Peer) ListenPeer() {

	var pListenAddr, err = net.ResolveTCPAddr("tcp4", peer.LocalAddr+":9093")
	CheckError(err)
	fmt.Println(pListenAddr.String())

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

	f, err := os.Open(peer.MyFolderPath+"/"+peer.MyFiles[msg.RootHash].Name)
	CheckError(err)
	defer f.Close()

	// Ovde sajlem File objekat koji downloader prima
	finfo, err := f.Stat()
	CheckError(err)

	path := peer.MyFolderPath+"/"+peer.MyFiles[msg.RootHash].Name
	tmpWriter.WriteFile(path, msg.ChunkNum, *peer.MyFiles[msg.RootHash].ChunkSize, finfo.Size())
}

func (peer Peer) connectToPeer(fname string, IP string, group *sync.WaitGroup, f *os.File, numOfPart int64, chunkStatuses []int, mutex *sync.Mutex, numOfDownloaded *int64) {
	fmt.Printf("[connectToPeer] About to dial: %+v\n", IP)

	rAddr, err := net.ResolveTCPAddr("tcp", IP+":9093")

	conn, err := net.DialTCP("tcp", nil, rAddr)
	CheckError(err)

	tmpReader := IO.Reader{conn}
	tmpWriter := IO.Writer{conn}

	// Hardkodovan root hash fajla koji hocu

	msg := MsgToNode{fname, numOfPart}

	msgForSend, err := json.Marshal(msg)
	CheckError(err)

	tmpWriter.Write(string(msgForSend))

	partBytes, size := tmpReader.ReadFile()

	//fmt.Printf("wg: %+v\n", group)

	// Ovde mora da se zakljuca fajl pre pisanja
	mutex.Lock()

	_, err = f.WriteAt(partBytes[:size], int64(numOfPart) * size)

	chunkStatuses[numOfPart] = 2

	*numOfDownloaded++

	mutex.Unlock()

	CheckError(err)

	group.Done()
}

