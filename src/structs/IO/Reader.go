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

func appendBuffer(dest []byte, src []byte, offset int64, length int64) []byte {
	//fmt.Println("offset:", offset, "len:", length, "src len:", len(src))
	var i int64
	j := 0
	for i = offset; i < length + offset; i++ {
		//fmt.Println("Lepim na indeks ", i)
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

	fmt.Println("In", r.Conn.RemoteAddr(), "psize:", pSize)

	for i := 0; ; i++ {
		n, err := r.Conn.Read(buff)

		CheckError(err)

		//fmt.Printf("[ReadFile] %+v-ti read Od %+v sam dobio bajtove: %+v\n", i, r.Conn.RemoteAddr(), n)

		finalBuff = appendBuffer(finalBuff, buff[:n], sum, int64(n))

		CheckError(err)

		sum += int64(n)
		if sum == int64(pSize) {
			fmt.Println("About to break in:",r.Conn.RemoteAddr(), "sum:", sum)
			break
		}
	}

	return finalBuff[:sum], sum
}
