package util

import "time"

const (
	Success = 0
	Fail    = 1

	MessCreateUserFailReason = 4

	ReadTimeout = time.Second * 5

	SendChanCnt    = 10
	RecvChanCnt    = 10
	ProcessChanCnt = 10

	MaxUsernameLen = 40

	ExitChatFlag = "###"
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
)

// client menu status
const (
	NoUser = iota
	Select
	ChatRoom
)
