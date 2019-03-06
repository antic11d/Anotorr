package Node

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
)

type Peer struct {

	Id string
	IP net.IP
	PrivateKey *rsa.PrivateKey
	RootHashes []string
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


