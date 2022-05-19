package process

import (
	"bufio"
	"client/util"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Process(recv, send chan []byte, ctx context.Context) {
	command := make(chan string, 0)
	var currMenu int = util.NoUser
	var iMessage util.IMessage
	go processInput(command)
	displayTip(currMenu)
	for {
		select {
		case <-ctx.Done():
			log.Printf("recv exit signal,exit")
			return
		case data := <-recv:
			funcCode := data[0]
			switch funcCode {
			case util.CreateUserAck:
				success := data[1]
				switch success {
				case util.Success:
					personId := binary.BigEndian.Uint16(data[2:4])
					iMessage.UserId = personId
					fmt.Printf("congratulation! %s register success \n", iMessage.Username)
					currMenu = util.Select
				case util.Fail:
					fmt.Printf("register fail,reason:%s\n", string(data[2:]))
					//fmt.Println("Please input another username again:")
				default:
				}
			case util.ChatListAck:
				splitChar := data[1]
				charList := strings.Split(string(data[2:]), string(splitChar))
				fmt.Println("*****Chat Room List*****")
				fmt.Println("*Chat id|Chat name*")
				for _, s := range charList {
					fmt.Println(s)
				}
				fmt.Println("******End******")
			//fmt.Printf("Please select one of the following:\n  0: View a list of chat rooms\n  1:Enter char room\n")
			case util.EnterChatAck:
				success := data[1]
				switch success {
				case util.Success:
					fmt.Printf("*****Welcome*****\n")
					fmt.Printf("***You could input '%s' to leave chat room***\n", util.ExitChatFlag)
					currMenu = util.ChatRoom
				case util.Fail:
					fmt.Println(string(data[2:]), "\nEnter chatRoom fail,please try it again")
				default:
				}
			case util.ExitChatAck:
				success := data[1]
				switch success {
				case util.Success:
					fmt.Println("*****Exit chatRoom*****")
					currMenu = util.Select
				case util.Fail:
					fmt.Printf("*****Exit chatRoom fail,%s*****\n", string(data[2:]))
				default:
				}
			case util.SendMessAck:
				// nothing to do
				//fmt.Printf("[YOU]:%s\n", iMessage.ChatMess)
				success := data[8]
				switch success {
				case util.Success:
				case util.Fail:
					fmt.Printf("send message fail,%s\n", string(data[9:]))
				default:
				}
			case util.CreateChatAck:
				success := data[1]
				switch success {
				case util.Success:
					fmt.Println("Create success.")
				case util.Fail:
					fmt.Printf("%s create Fail,%s\n", iMessage.ChatRoomName, string(data[2:]))
				default:
				}
			case util.OtherPersonMess:
				nameLength := binary.BigEndian.Uint16(data[1:3])
				person := string(data[3 : 3+nameLength])
				mess := string(data[3+nameLength:])
				fmt.Printf("[%s %s]:%s\n", time.Now().Format("2006-01-02 15:04:05"), person, mess)
			case util.OtherPersonEnterChat:
				person := string(data[1:])
				fmt.Printf("------%s Welcome <%s> to join the group chat!-----\n", time.Now().Format("2006-01-02 15:04:05"), person)
			default:
				fmt.Println("*****Unknown error*****")
			}
			displayTip(currMenu)
		case cmd := <-command:
			switch currMenu {
			case util.NoUser:
				iMessage.Username = cmd
				send <- iMessage.PackCreateUser()
			case util.Select:
				switch cmd {
				case "0":
					send <- iMessage.PackChatList()
				case "1":
					fmt.Println("Please input chat id:")
					cmd = <-command
					chatId, err := strconv.Atoi(cmd)
					if err != nil {
						fmt.Println(cmd, "isn't a valid chatId")
						displayTip(currMenu)
						continue
					}
					iMessage.ChatId = uint16(chatId)
					send <- iMessage.PackEnterChat()
				case "2":
					fmt.Println("Please input chat room name:")
					cmd = <-command
					iMessage.ChatRoomName = cmd
					send <- iMessage.PackCreateChat()
				default:
					fmt.Println("input error,please input again!!!")
				}
			case util.ChatRoom:
				if cmd == util.ExitChatFlag {
					send <- iMessage.PackExitChat()
					continue
				}
				iMessage.ChatMess = cmd
				send <- iMessage.PackChatMess()
			}
		}
	}
}

func processInput(command chan string) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		command <- input
	}
}

func displayTip(currMenu int) {
	switch currMenu {
	case util.NoUser:
		fmt.Println("Please input your username:")
	case util.Select:
		fmt.Printf("Please select one of the following:\n  0:View a list of chat rooms\n  1:Enter a chat room\n" +
			"  2:Create a chat room\n")
	case util.ChatRoom:
	default:
	}
}
