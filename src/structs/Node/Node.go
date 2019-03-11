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
	Id string
	IP net.IP
	PrivateKey *rsa.PrivateKey
	RootHashes []string
	TrackerListenAddr *net.TCPAddr
	PeerListenAddr *net.TCPAddr
	ListenerTracker *net.TCPListener
	ListenerPeer *net.TCPListener
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getMyIp() (net.IP) {

	//Ovde ce da se implementira UPNP

	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
			return ipv4
		}
	}

	return nil
}

func InitializeNode() (p *Peer){

	pk, err := rsa.GenerateKey(rand.Reader, 2048)

	CheckError(err)

	p = &Peer{Id:"idPrvi", PrivateKey:pk, IP:getMyIp()}

	return p

}

// Hardcoded portovi 9091, 9092

// Cekam da mi se javi treker da mi kaze da neko hoce da skida fajl koji ja potencijalno imam
func (peer Peer) ListenTracker() {
	fmt.Printf("Hello from listentracker!!!")
	var tListenAddr, err = net.ResolveTCPAddr("tcp4", ":9091")
	CheckError(err)

	peer.ListenerTracker, err = net.ListenTCP("tcp", tListenAddr)
	CheckError(err)

	for  {
		conn, err := peer.ListenerTracker.Accept()
		fmt.Printf("Got a call from tracker!")
		if err != nil {
			fmt.Println("Error while accepting connection from tracker, continuing...")
			continue
		}

		go handleTracker(conn)
	}
}

func handleTracker(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Hello from handletracker, conn: %+v\n", conn)
	var reader = IO.Reader{conn}
	var writer = IO.Writer{conn}

	msg := reader.Read()

	var wrappedRequest Requests.WrappedRequest

	err := json.Unmarshal([]byte(msg), &wrappedRequest)
	CheckError(err)

	fmt.Printf("Dobio objekat: %+v\n", wrappedRequest)

	// Ovde treba da kazem imam taj fajl, i da vratim svoj IP trekeru
	writer.Write("192.168.0.28")
}


