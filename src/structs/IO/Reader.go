package IO

import (
	"net"
	"strconv"
)

type Reader struct {
	Conn *net.TCPConn
}

func (r Reader) Read() string  {
	recvBuff := make([]byte, 2048)
	bytesRead,err := r.Conn.Read(recvBuff)
	if err != nil {
		panic(err)
	}

	return string(recvBuff[:bytesRead])
}

func appendBuffer(dest []byte, src []byte, offset int64, length int64) []byte {
	var i int64
	j := 0
	for i = offset; i < length + offset; i++ {
		dest[i] = src[j]
		j++
	}

	return dest
}

func (r Reader) ReadFile() ([]byte, int64) {
	partSize := r.Read()
	r.Conn.Write([]byte("ok"))

	pSize, err := strconv.Atoi(partSize)
	CheckError(err)

	buff := make([]byte, 32)
	finalBuff := make([]byte, pSize*2)

	CheckError(err)

	var sum int64 = 0

	for i := 0; ; i++ {
		n, err := r.Conn.Read(buff)

		CheckError(err)

		finalBuff = appendBuffer(finalBuff, buff[:n], sum, int64(n))

		CheckError(err)

		sum += int64(n)
		if sum == int64(pSize) {
			break
		}
	}

	return finalBuff[:sum], sum
}
