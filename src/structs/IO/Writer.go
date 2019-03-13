package IO

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Writer struct {
	Conn net.Conn
}

func (w Writer) Write(msg string) {
	_,err := w.Conn.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}
func (w Writer) WriteFile(filename string, offset int64, partSize int64){

	f, err := os.Open("misc/txtFile.txt")
	defer f.Close()
	CheckError(err);
	_,err = f.Seek(partSize*offset,0)
	CheckError(err)

	tmpPartSize := strconv.FormatInt(partSize,10)
	bufferSize, err := strconv.Atoi(tmpPartSize)
	CheckError(err)
	buffer := make([]byte, bufferSize)
	bytesRead, err := f.Read(buffer)
	CheckError(err)

	fmt.Println(string(tmpPartSize))

	w.Write(string(tmpPartSize))

	tmpBuffer := make([]byte, 3)
	_, err = w.Conn.Read(tmpBuffer)
	CheckError(err)

	_, err = w.Conn.Write(buffer[:bytesRead])
	CheckError(err)
}
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
