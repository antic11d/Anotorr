package Tracker

import "net"

type Tracker struct {
	Addr *net.TCPAddr
	ListOfItems []string
}
