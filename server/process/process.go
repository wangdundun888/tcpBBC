package process

import (
	"context"
	"encoding/binary"
	"log"
	"math"
	"net"
	"server/util"
	"strconv"
)

type person struct {
	conn     net.Conn
	pid      uint16
	username string
}

func Process(recv, send chan util.IMessage, ctx context.Context) {
	// manage all tcp conns
	connMap := make(map[string]person, 100)
	// manage all users
	userExistMap := make(map[string]struct{}, 100)
	// manage all chat rooms
	chatExistMap := make(map[string]uint16, 100)
	// manage all in the same chat room users
	onlineChatMap := make(map[uint16][]person, 100)
	for {
		select {
		case <-ctx.Done():
			log.Printf("recv exit signal,exit")
			return
		case mess := <-recv:
			data := mess.Data
			funcCode := data[0]
			conn := mess.Conn
			switch funcCode {
			case util.Connect:
				addr := mess.Conn.RemoteAddr().String()
				p := person{
					conn: mess.Conn,
				}
				connMap[addr] = p
			case util.DisConnect:
				addr := mess.Conn.RemoteAddr().String()
				p, exist := connMap[addr]
				if exist {
					delete(connMap, addr)
					delete(userExistMap, p.username)
					for cid, people := range onlineChatMap {
						var newPeople []person
						for _, p2 := range people {
							if p2.conn.RemoteAddr().String() != addr {
								newPeople = append(newPeople, p2)
							}
						}
						onlineChatMap[cid] = newPeople
					}
				}
			case util.CreateUser:
				username := string(data[1:])
				_, exist := userExistMap[username]
				if exist {
					resp := util.PackCreateUserResp(util.Fail, 0, "username exist")
					mess.Data = resp
					send <- mess
					continue
				}
				pid := len(userExistMap)
				if pid >= math.MaxUint16 {
					resp := util.PackCreateUserResp(util.Fail, 0, "service busy please try it again latter")
					mess.Data = resp
					send <- mess
					continue
				}
				userExistMap[username] = struct{}{}
				p := connMap[conn.RemoteAddr().String()]
				p.pid = uint16(pid)
				p.username = username
				connMap[conn.RemoteAddr().String()] = p
				resp := util.PackCreateUserResp(util.Success, uint16(pid), "")
				mess.Data = resp
				send <- mess
			case util.CreateChat:
				chatRoomName := string(data[1:])
				_, exist := chatExistMap[chatRoomName]
				if exist {
					resp := util.PackCreateChatResp(util.Fail, 0, "chat room exist")
					mess.Data = resp
					send <- mess
					continue
				}
				cid := uint16(len(chatExistMap))
				chatExistMap[chatRoomName] = cid
				resp := util.PackCreateChatResp(util.Success, cid, "")
				mess.Data = resp
				send <- mess
			case util.ChatList:
				var chatList []string
				for name, cid := range chatExistMap {
					chatList = append(chatList, strconv.Itoa(int(cid))+" "+name)
				}
				resp := util.PackChatListResp(chatList, util.SplitChar)
				mess.Data = resp
				send <- mess
			case util.EnterChat:
				var exist bool
				cid := binary.BigEndian.Uint16(data[1:3])
				for _, id := range chatExistMap {
					if cid == id {
						exist = true
						break
					}
				}
				if !exist {
					resp := util.PackEnterChatResp(util.Fail, 0, "chat room no exist")
					mess.Data = resp
					send <- mess
					continue
				}
				personSlice, exist := onlineChatMap[cid]
				p := connMap[mess.Conn.RemoteAddr().String()]
				// notice other person who in the same chat room
				for _, p2 := range personSlice {
					resp := util.PackOtherEnterChatResp(p.username)
					iMess := util.IMessage{
						Conn: p2.conn,
						Data: resp,
					}
					send <- iMess
				}
				// enter ack
				if exist {
					personSlice = append(personSlice, p)
					onlineChatMap[cid] = personSlice
				} else {
					onlineChatMap[cid] = []person{p}
				}
				resp := util.PackEnterChatResp(util.Success, cid, "")
				mess.Data = resp
				send <- mess
			case util.ExitChat:
				cid := binary.BigEndian.Uint16(data[1:])
				oldPersonSlice := onlineChatMap[cid]
				var newPersonSlice []person
				for _, p := range oldPersonSlice {
					if p.conn.RemoteAddr().String() != mess.Conn.RemoteAddr().String() {
						newPersonSlice = append(newPersonSlice, p)
					}
				}
				onlineChatMap[cid] = newPersonSlice
				resp := util.PackExitChatResp(util.Success, cid, "")
				mess.Data = resp
				send <- mess
			case util.SendMess:
				cid := binary.BigEndian.Uint16(data[1:3])
				pid := binary.BigEndian.Uint16(data[3:5])
				personMess := string(data[8:])
				// mess
				resp := util.PackSendMessAckResp(cid, pid, util.Success, "")
				mess.Data = resp
				send <- mess
				// send to others who in the same chat room
				personSlice, _ := onlineChatMap[cid]
				p := connMap[mess.Conn.RemoteAddr().String()]
				resp = util.PackSendOtherMessResp(p.username, personMess)
				for _, p2 := range personSlice {
					if p.conn.RemoteAddr().String() != p2.conn.RemoteAddr().String() {
						iMess := util.IMessage{
							Conn: p2.conn,
							Data: resp,
						}
						send <- iMess
					}
				}
			default:
			}
		}
	}
}
