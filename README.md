# goTorr

### Project summary

* The idea is to implement peer-to-peer torrent protocol. Every node (peer) in network can be both uploader and downloader.

* Tracker server will contain 'blocks' with root hash that will be used for proof of consistance and existance.  
For that cause, Merkle Tree data structure will be used. It will also contain queue with requests for downloading and uploading files. Each file will be downloaded in fixed number of chunks, because Merkle Tree works with fixed number of hashes.  
Tracker will also contain 'blacklist' of nodes who had corrupted file at some point in time.   
In future, tracker will be replaced with Ethereum blockchain, because we want decentralization of network.

* UPnP protocol is used for connecting nodes in network.

***
### Language in use:  
* __Go__
***
### Authors  
* Antić Dimitrije 128/2016  
* Golubović Stefan 135/2016  
* Novaković Andrija 68/2016  
