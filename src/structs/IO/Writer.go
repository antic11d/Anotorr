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
func (w Writer) WriteFile(filename string, offset int64, partSize int64){
	f, err := os.Open(filename)
	defer f.Close()
	CheckError(err);
	_,err = f.Seek(partSize*offset,0)
	CheckError(err)

	tmpPartSize := strconv.FormatInt(partSize,10)

	CheckError(err)

	// Saljem downloaderu velicinu parta
	w.Write(string(tmpPartSize))

	// Ovde citam OK da je stigla velicina fajla
	tmpBuffer := make([]byte, 3)
	_, err = w.Conn.Read(tmpBuffer)
	fmt.Println()
	CheckError(err)

	fmt.Printf("[WriteFile] tmpbuffer pre iscitavanja fajla: %+v\n",tmpBuffer)

	//msg := make([]byte, 5)
	tmpBuffer = make([]byte, 256)
	var bytesSent int64 = 0

	for bytesSent < partSize {
		if partSize - bytesSent < 256 {
			tmpBuffer = make([]byte,partSize - bytesSent)
		}
		bytesRead, err := f.Read(tmpBuffer)

		n, err := w.Conn.Write(tmpBuffer[:bytesRead])
		CheckError(err)

		//_, err = w.Conn.Read(msg)

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