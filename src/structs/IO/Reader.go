package IO

import (
	"fmt"
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

func appendBuffer(dest []byte, src []byte, offset int, length int) []byte {
	finLen := length + offset
	for i := offset; i < finLen; i++ {
		dest[i] = src[i-offset]
	}

	return dest
}

func (r Reader) ReadFile() ([]byte, int64) {
	partSize := r.Read()
	r.Conn.Write([]byte("ok"))

	pSize, err := strconv.Atoi(partSize)
	CheckError(err)

	buff := make([]byte, 1024)
	finalBuff := make([]byte, pSize)

	CheckError(err)

	var sum int64 = 0

	fmt.Println(pSize)

	for i := 0; ; i++ {
		n, err := r.Conn.Read(buff)
		CheckError(err)

		//fmt.Printf("[ReadFile] %+v-ti read Od %+v sam dobio bajtove: %+v\n", i, r.Conn.RemoteAddr(), n)

		finalBuff = appendBuffer(finalBuff, buff[:n], int(sum), n)

		CheckError(err)

		sum += int64(n)
		if sum == int64(pSize) {
			fmt.Println("About to break:", sum)
			break
		}

		CheckError(err)
	}

	return finalBuff, sum
}
