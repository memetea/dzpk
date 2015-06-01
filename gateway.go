package dzpk

import (
	"bufio"
	"net"
)

type PlayerConnection struct {
	Conn      net.Conn
	WriteChan chan<- []byte
	ReadChan  <-chan []byte
}

func handleClientConn(conn net.Conn) {
	br := bufio.NewReader(conn)
	version, _ := br.ReadByte()
	if version != 1 {
		//conn.Write(b)
		conn.Close()
	}

}

func StartGateway() {
	ln, err := net.Listen("tcp", ":8282")
	if err != nil {
		//log.Println("tcp listen on :8282 failed:")
		syslog.Error("listen on :8282 error:%v", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			syslog.Warning("Accept client err:%v", err)
			continue
		}

		go handleClientConn(conn)
	}
}
