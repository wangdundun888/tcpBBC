package recv

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"server/util"
)

func Recv(tcpConn net.Conn, process chan util.IMessage, ctx context.Context) {
	defer func() {
		data := []byte{util.DisConnect}
		iMess := util.IMessage{
			Conn: tcpConn,
			Data: data,
		}
		process <- iMess
		fmt.Println(tcpConn.RemoteAddr(), "exit...")
	}()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Recv goroutine exit\n")
			return
		default:
		}
		lengthBytes := make([]byte, 2, 2)
		_, err := tcpConn.Read(lengthBytes)
		if err != nil {
			//log.Printf("read head error %s\n", err.Error())
			return
		}
		length := binary.BigEndian.Uint16(lengthBytes)
		data := make([]byte, length, length)
		n, err := tcpConn.Read(data)
		if err != nil {
			log.Printf("read data error %s\n", err.Error())
			return
		}
		if n < int(length) {
			log.Printf("read data error,read %v < real %v\n", n, length)
			continue
		}
		iMess := util.IMessage{
			Conn: tcpConn,
			Data: data[:n],
		}
		process <- iMess
	}
}
