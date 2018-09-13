package main

import (
	"github.com/name5566/leaf/cluster"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/console"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"
	"os"
	"os/signal"
	"leaf_server/oss"
)

func Run(mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	osslogger, err := oss.New("udp", conf.OssLogAddr, "test")
	if nil != err {
		log.Error("init oss logger failed:", err.Error())
		return
	}
	oss.Export(osslogger)
	defer oss.Close()

	
	log.Release("starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
