package Tracker

import
(	"net"
	"../File"
)

type Tracker struct {
	Addr *net.TCPAddr
	Map map[string]File.File
}
