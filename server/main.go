package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"server/process"
	"server/recv"
	"server/send"
	"server/util"
	"time"
)

func main() {
	cmd := NewCmd()
	tcpListener, err := net.Listen("tcp", ":"+cmd.Port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer tcpListener.Close()
	processChan := make(chan util.IMessage, util.ProcessChanCnt)
	sendChan := make(chan util.IMessage, util.SendChanCnt)
	bg := context.Background()
	ctx, cancel := context.WithCancel(bg)
	go send.Send(sendChan, ctx)
	go process.Process(processChan, sendChan, ctx)
	defer func() {
		cancel()
		time.Sleep(time.Second)
		log.Printf("exit\n")
	}()
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		//  new connection
		data := []byte{util.Connect}
		iMess := util.IMessage{
			Conn: conn,
			Data: data,
		}
		processChan <- iMess
		go recv.Recv(conn, processChan, ctx)
	}

}
