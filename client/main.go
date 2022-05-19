package main

import (
	"client/process"
	"client/recv"
	"client/send"
	"client/util"
	"context"
	"log"
	"net"
	"time"
)

func main() {
	cmd := NewCmd()
	tcpConn, err := net.Dial("tcp", cmd.Ip)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return
	}
	defer tcpConn.Close()
	bg := context.Background()
	ctx, cancel := context.WithCancel(bg)
	sendChan := make(chan []byte, util.SendChanCnt)
	processChan := make(chan []byte, util.ProcessChanCnt)
	go recv.Recv(tcpConn, processChan, ctx)
	go send.Send(tcpConn, sendChan, ctx)
	go process.Process(processChan, sendChan, ctx)
	<-ctx.Done()
	defer func() {
		cancel()
		time.Sleep(time.Second)
		log.Printf("exit\n")
	}()
}
