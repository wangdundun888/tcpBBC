package util

import (
	"encoding/binary"
)

type Message interface{}
type IMessage struct {
	Username string
	UserId   uint16
	Length   uint16
	MessType byte
	Data     []byte

	ChatRoomName string
	ChatId       uint16
	ChatMess     string
}

func (im *IMessage) PackCreateUser() []byte {
	var usernameBytes []byte
	usernameBytes = append(usernameBytes, []byte(im.Username)...)
	data := make([]byte, 2, 2)
	// length
	binary.BigEndian.PutUint16(data, uint16(len(usernameBytes))+1)
	// message type
	data = append(data, CreateUser)
	// data
	data = append(data, usernameBytes...)
	return data
}

func (im *IMessage) PackChatList() []byte {
	var dataBytes []byte
	dataBytes = append(dataBytes, ChatList)
	lengthBytes := make([]byte, 2, 2)
	length := uint16(1)
	binary.BigEndian.PutUint16(lengthBytes, length)
	dataBytes = append(lengthBytes, dataBytes...)
	return dataBytes
}

func (im *IMessage) PackEnterChat() []byte {
	var dataBytes []byte
	dataBytes = append(dataBytes, EnterChat)
	dataBytes = append(dataBytes, []byte{0, 0}...)
	binary.BigEndian.PutUint16(dataBytes[1:], im.ChatId)
	lengthBytes := make([]byte, 2, 2)
	length := uint16(len(dataBytes))
	binary.BigEndian.PutUint16(lengthBytes, length)
	dataBytes = append(lengthBytes, dataBytes...)
	return dataBytes
}

func (im *IMessage) PackExitChat() []byte {
	var dataBytes []byte
	dataBytes = append(dataBytes, ExitChat)
	dataBytes = append(dataBytes, []byte{0, 0}...)
	binary.BigEndian.PutUint16(dataBytes[1:], im.ChatId)
	lengthBytes := make([]byte, 2, 2)
	length := uint16(len(dataBytes))
	binary.BigEndian.PutUint16(lengthBytes, length)
	dataBytes = append(lengthBytes, dataBytes...)

	return dataBytes
}

func (im *IMessage) PackChatMess() []byte {
	var dataBytes []byte
	dataBytes = append(dataBytes, SendMess)
	dataBytes = append(dataBytes, []byte{0, 0}...)
	binary.BigEndian.PutUint16(dataBytes[1:], im.ChatId)
	dataBytes = append(dataBytes, []byte{0, 0}...)
	binary.BigEndian.PutUint16(dataBytes[3:], im.UserId)
	dataBytes = append(dataBytes, []byte{0, 0, 0}...)
	dataBytes = append(dataBytes, []byte(im.ChatMess)...)
	lengthBytes := make([]byte, 2, 2)
	length := uint16(len(dataBytes))
	binary.BigEndian.PutUint16(lengthBytes, length)
	dataBytes = append(lengthBytes, dataBytes...)

	return dataBytes
}

func (im *IMessage) PackCreateChat() []byte {
	var dataBytes []byte
	dataBytes = append(dataBytes, CreateChat)
	dataBytes = append(dataBytes, []byte(im.ChatRoomName)...)
	lengthBytes := make([]byte, 2, 2)
	length := uint16(len(dataBytes))
	binary.BigEndian.PutUint16(lengthBytes, length)
	dataBytes = append(lengthBytes, dataBytes...)

	return dataBytes
}
