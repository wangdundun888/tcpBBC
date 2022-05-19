package send

import (
	"context"
	"log"
	"net"
)

func Send(tcpConn net.Conn, send chan []byte, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Send goroutine exit\n")
			return
		case mess := <-send:
			go write(tcpConn, mess)
		}
	}
}

func write(tcpConn net.Conn, mess []byte) {
	_, err := tcpConn.Write(mess)
	if err != nil {
		log.Printf("send mess error:%s\n", err.Error())
	}
}
