package util

import "time"

const (
	Success = 0
	Fail    = 1

	ReadTimeout = time.Second * 5

	MaxUsernameLen = 40

	ProcessChanCnt = 100
	RecvChanCnt    = 100
	SendChanCnt    = 100

	SplitChar = "#"
)

// message type
const (
	ChatList = iota
	ChatListAck
	EnterChat
	EnterChatAck
	ExitChat
	ExitChatAck
	CreateChat
	CreateChatAck
	SendMess
	SendMessAck
	CreateUser
	CreateUserAck
	OtherPersonMess
	OtherPersonEnterChat

	Connect
	DisConnect
)
