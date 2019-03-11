package Node

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"../IO"
	"../Requests"
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
	return &Peer{ID:"idPrvi", PrivateKey:pk, IP: getMyIP()}

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
	fmt.Printf("[handleTracker] lista: %+v\n", wrappedRequest.Value.CryptedIPs.Len())

	// Ovde treba da prodjem kroz svoji fajl sistem i da vidim da li imam taj fajl, ako imam onda vratim svoj IP trekeru
	tmpWriter.Write("10.0.151.148")
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

	fmt.Printf("[RequestDownload] Treba da se javim svima iz liste: %+v i duzine %+v\n", completedReq.Value.CryptedIPs, completedReq.Value.CryptedIPs.Len())
}


