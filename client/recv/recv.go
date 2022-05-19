package recv

import (
	"client/util"
	"context"
	"encoding/binary"
	"log"
	"net"
	"time"
)

func Recv(tcpConn net.Conn, process chan []byte, ctx context.Context) {
	if tcpConn == nil {
		log.Printf("tcpConn is nil\n")
		return
	}
	for {
		select {
		case <-ctx.Done():
			log.Printf("Recv goroutine exit\n")
			return
		default:
		}
		tcpConn.SetDeadline(time.Now().Add(util.ReadTimeout))
		lengthBytes := make([]byte, 2, 2)
		_, err := tcpConn.Read(lengthBytes)
		if err != nil {
			//log.Printf("read head error %s\n", err.Error())
			continue
		}
		length := binary.BigEndian.Uint16(lengthBytes)
		data := make([]byte, length, length)
		n, err := tcpConn.Read(data)
		if err != nil {
			log.Printf("read data error %s\n", err.Error())
			continue
		}
		if n < int(length) {
			log.Printf("read data error,read %v < real %v\n", n, length)
			continue
		}
		process <- data[:n]
	}
}
