package Requests

import (
	"container/list"
	"crypto/rsa"
)

type DownloadRequestKey struct {
	RootHash  string
	PublicKey *rsa.PublicKey
}

type DownloadRequest struct {
	CryptedIPs *list.List
	Served int // 0 - lista seedera je prazna, ako imas fajl kriptuj i vrati mi (treker je ja)
				// 1 - saljem ti listu ljudi od kojih mozes da skidas (imas ih u CryptedIPs (treker je ja)
}

type WrappedRequest struct {
	Key DownloadRequestKey
	Value DownloadRequest
}
