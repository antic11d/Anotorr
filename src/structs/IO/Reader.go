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

func (r Reader) ReadFile() ([]byte, int64) {
	partSize := r.Read()
	r.Conn.Write([]byte("ok"))

	pSize, err := strconv.Atoi(partSize)
	CheckError(err)

	buff := make([]byte, 256)
	finalBuff := make([]byte, 0)

	CheckError(err)

	var sum int64 = 0

	fmt.Println(pSize)

	for i := 0; ; i++ {
		n, err := r.Conn.Read(buff)
		CheckError(err)

		_, err = r.Conn.Write([]byte("next"))
		CheckError(err)

		fmt.Printf("[ReadFile] %+v-ti read Od %+v sam dobio bajtove: %+v\n", i, r.Conn.RemoteAddr(), n)

		finalBuff = append(finalBuff, buff[:n]...)

		sum += int64(n)
		if sum == int64(pSize) {
			fmt.Println("About to break:", sum)
			break
		}

		CheckError(err)
	}


	//rBuff := finalBuff[:sum]
	return finalBuff, sum
}
