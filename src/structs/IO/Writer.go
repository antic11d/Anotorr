package IO

import (
	"net"
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
