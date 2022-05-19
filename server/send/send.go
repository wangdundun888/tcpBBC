package send

import (
	"context"
	"log"
	"net"
	"server/util"
)

func Send(send chan util.IMessage, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Send goroutine exit\n")
			return
		case mess := <-send:
			go write(mess.Conn, mess.Data)
		}
	}
}

func write(tcpConn net.Conn, mess []byte) {
	_, err := tcpConn.Write(mess)
	if err != nil {
		log.Printf("send mess error:%s\n", err.Error())
	}
}
