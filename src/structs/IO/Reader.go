package IO

import (
	"fmt"
	"net"
	"strconv"
)

type Reader struct {
	Conn net.Conn
}

func (r Reader) Read() string  {
	recvBuff := make([]byte, 2048)
	bytesRead,err := r.Conn.Read(recvBuff)
	if err != nil {
		panic(err)
	}

	return string(recvBuff[:bytesRead])
}

func (r Reader) ReadFile() ([]byte, int) {
	partSize := r.Read()
	r.Conn.Write([]byte("ok"))

	pSize, err := strconv.Atoi(partSize)
	CheckError(err)

	buff := make([]byte, pSize)
	finalBuff := make([]byte, pSize)

	CheckError(err)

	var sum int64 = 0

	fmt.Println(pSize)

	for i := 0; ; i++ {
		n, err := r.Conn.Read(buff)
		CheckError(err)

		fmt.Printf("[ReadFile] %+v-ti read Od %+v sam dobio bajtove: %+v\n", i, r.Conn.RemoteAddr(), n)

		sum += int64(n)
		if sum == int64(pSize) {
			break
		}

		finalBuff = append(finalBuff, buff...)

		CheckError(err)
	}

	fmt.Println("[ReadFile] Done reading!")

	return finalBuff, len(finalBuff)
}
