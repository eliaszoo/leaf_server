package main

import (
	"leaf_server/conf"
	"leaf_server/game"
	"leaf_server/gate"
	"leaf_server/login"

	lconf "github.com/name5566/leaf/conf"
)


func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	lconf.OssLogAddr  = conf.Server.OssLogAddr

	Run(
		gate.Module,
		game.Module,
		login.Module,
	)
}
