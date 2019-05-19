package Requests

import (
	"crypto/rsa"
)

type DownloadRequestKey struct {
	RootHash  string
	PublicKey *rsa.PublicKey
}

type DownloadRequest struct {
	CryptedIPs Matrix
	Served *int
}

type WrappedRequest struct {
	Key *DownloadRequestKey
	Value *DownloadRequest
}

type Matrix struct {
	Arr [][]byte
}
