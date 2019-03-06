package Requests

import "crypto/rsa"

type DownloadRequestKey struct {
	RootHash  string
	PublicKey *rsa.PublicKey
}

type DownloadRequest struct {
	CryptedIPs []string
	Served bool
}
