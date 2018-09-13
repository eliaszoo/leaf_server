package gate

import (
	"leaf_server/game"
	"leaf_server/login"
	"leaf_server/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.LoginMsg{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.GameMsg{}, game.ChanRPC)
}
