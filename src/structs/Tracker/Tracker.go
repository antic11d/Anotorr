package Tracker

import (
	"../File"
	"../Requests"
	"net"
)

type Tracker struct {
	Addr *net.TCPAddr
	Map map[string]File.File
	DownloadRequests map[Requests.DownloadRequestKey] Requests.DownloadRequest
}
