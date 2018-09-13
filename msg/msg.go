package msg

import (
	"leaf_server/msg/processor"
)

var Processor = processor.NewProcessor()

func init() {
	Processor.Register(&LoginMsg{})
	Processor.Register(&GameMsg{})
	Processor.Register(&RetMsg{})
}

type LoginMsg struct {
	Cmd string      `json:"cmd"`
	Req interface{} `json:"req"`
}

type GameMsg struct {
	Cmd string		`json:"cmd"`
	Req interface{}	`json:"req"`
}

type RetMsg struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errormsg"`
	Cmd      string      `json:"cmd"`
	Ans      interface{} `json:"data"`
}

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginAns struct {
	UserCheck string `json:"usercheck"`
}

type TestReq struct {
}

type TestAns struct {
}