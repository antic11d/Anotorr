# goTorr

### Project summary

* The idea is to implement peer-to-peer torrent protocol. 

* Tracker server works like a signal server, when someone wants to download a file, tracker reacts on posted request and contacts all peers in network, so they can seed that file.  
For consistence, Merkle Tree data structure will be used. Tracker will also contain queue with requests for downloading and uploading files. Each file will be downloaded in fixed number of chunks.  
Tracker will discredit nodes that have "changed" file. Important aspect is that tracker does not contains IPs of every node that contains each file which is the main difference with other torrent protocols. 
<br>In future, signal server will be replaced with Ethereum blockchain, because we want decentralization of network.
* Every node (peer) in network can be both seeder and downloader.
* uPnP protocol is used for connecting nodes in network, where user explicitly forward ports.
* Known issue is that uPnP protocol isn't supported on all router devices, also mobile network operator does not support uPnP so our app is (for now) available only for pc.
***
### Platforms:
* Our project is being developed and tested on Linux operating systems, but there is no reason it would not work on Windows.  
***
### Language in use:  
* __Go__
***
### Authors  
* Antić Dimitrije 128/2016  
* Golubović Stefan 135/2016  
* Novaković Andrija 68/2016  
