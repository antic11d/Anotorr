package IO

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Writer struct {
	Conn *net.TCPConn
}

func (w Writer) Write(msg string) {
	_,err := w.Conn.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}

func calcSize(chunkNum int64, chunkSize int64, fileSize int64) int64 {
	if chunkNum * chunkSize + chunkSize > fileSize {
		return fileSize % chunkSize
	}

	return chunkSize
}

func (w Writer) WriteFile(filename string, chunkNum int64, chunkSize int64, fileSize int64){
	f, err := os.Open("misc/"+filename)
	defer f.Close()
	CheckError(err)

	sizeForSending := calcSize(chunkNum, chunkSize, fileSize)

	_, err = f.Seek(chunkSize*chunkNum,0)
	CheckError(err)

	// ne chunksize nego ono sto mi sracunamo
	tmpPartSize := strconv.FormatInt(sizeForSending,10)

	CheckError(err)

	// Saljem downloaderu velicinu parta
	w.Write(string(tmpPartSize))

	// Ovde citam OK da je stigla velicina fajla
	tmpBuffer := make([]byte, 3)
	_, err = w.Conn.Read(tmpBuffer)
	CheckError(err)

	tmpBuffer = make([]byte, 256)
	var bytesSent int64 = 0

	for bytesSent < sizeForSending {
		if sizeForSending - bytesSent < 256 {
			tmpBuffer = make([]byte, sizeForSending - bytesSent)
		}
		bytesRead, err := f.Read(tmpBuffer)

		n, err := w.Conn.Write(tmpBuffer[:bytesRead])
		CheckError(err)

		bytesSent += int64(n)

		fmt.Println("Poslao : ", bytesSent)
	}

	CheckError(err)
}
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}