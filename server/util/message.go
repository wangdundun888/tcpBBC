package util

import (
	"encoding/binary"
	"net"
	"strings"
)

type Message interface{}

type IMessage struct {
	Conn net.Conn
	Data []byte
}

func PackCreateUserResp(status byte, pid uint16, failReason string) []byte {
	var resp []byte
	resp = append(resp, CreateUserAck)
	resp = append(resp, status)
	if status == Success {
		pidBytes := make([]byte, 2, 2)
		binary.BigEndian.PutUint16(pidBytes, pid)
		resp = append(resp, pidBytes...)
	} else {
		resp = append(resp, []byte(failReason)...)
	}
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}

func PackCreateChatResp(status byte, cid uint16, failReason string) []byte {
	var resp []byte
	resp = append(resp, CreateChatAck)
	resp = append(resp, status)
	if status == Success {
		cidBytes := make([]byte, 2, 2)
		binary.BigEndian.PutUint16(cidBytes, cid)
		resp = append(resp, cidBytes...)
	} else {
		resp = append(resp, []byte(failReason)...)
	}
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}

func PackChatListResp(chatList []string, sep string) []byte {
	var resp []byte
	resp = append(resp, ChatListAck)
	resp = append(resp, []byte(sep)...)
	all := strings.Join(chatList, sep)
	resp = append(resp, []byte(all)...)
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}

func PackEnterChatResp(status byte, cid uint16, failReason string) []byte {
	var resp []byte
	resp = append(resp, EnterChatAck)
	resp = append(resp, status)
	if status == Success {
		cidBytes := make([]byte, 2, 2)
		binary.BigEndian.PutUint16(cidBytes, cid)
		resp = append(resp, cidBytes...)
	} else {
		resp = append(resp, []byte(failReason)...)
	}
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}

func PackExitChatResp(status byte, cid uint16, failReason string) []byte {
	var resp []byte
	resp = append(resp, ExitChatAck)
	resp = append(resp, status)
	if status == Success {
		cidBytes := make([]byte, 2, 2)
		binary.BigEndian.PutUint16(cidBytes, cid)
		resp = append(resp, cidBytes...)
	} else {
		resp = append(resp, []byte(failReason)...)
	}
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}
func PackOtherEnterChatResp(username string) []byte {
	var resp []byte
	resp = append(resp, OtherPersonEnterChat)
	resp = append(resp, []byte(username)...)
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}

func PackSendMessAckResp(cid, pid uint16, status byte, failReason string) []byte {
	var resp []byte
	resp = append(resp, SendMessAck)
	cidBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(cidBytes, cid)
	resp = append(resp, cidBytes...)
	pidBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(pidBytes, pid)
	resp = append(resp, pidBytes...)
	resp = append(resp, []byte{0, 0, 0}...)
	resp = append(resp, status)
	if status == Fail {
		resp = append(resp, []byte(failReason)...)
	}
	length := uint16(len(resp))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}

func PackSendOtherMessResp(username, mess string) []byte {
	var resp []byte
	resp = append(resp, OtherPersonMess)
	// username length
	length := uint16(len(username))
	lengthBytes := make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(resp, lengthBytes...)
	resp = append(resp, []byte(username)...)
	resp = append(resp, []byte(mess)...)
	// message length
	length = uint16(len(resp))
	lengthBytes = make([]byte, 2, 2)
	binary.BigEndian.PutUint16(lengthBytes, length)
	resp = append(lengthBytes, resp...)
	return resp
}
