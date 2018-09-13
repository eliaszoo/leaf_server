package internal

import (
	"leaf_server/base"
	"leaf_server/db"
	"leaf_server/conf"
	
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	mgodb    = new(db.Mongodb)
	PlayerManager = NewPlayerManager()
	TimerManager  = NewTimerManager()
	uidbuilder 	  = new(UidBuilder)
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	var err error
	mgodb, err = db.Dial(conf.Server.MgodbAddr, conf.Server.GameMgoConnNum, skeleton)
	if nil == mgodb {
		log.Error("dial mongodb failed:", conf.Server.MgodbAddr, " ", err.Error())
		return
	}

	mgodb.EnsureUniqueIndex(base.DBNAME, base.PLAYERSET, []string{"uid"})

	uidbuilder.Init()
}

func (m *Module) OnDestroy() {
	PlayerManager.Close()
	mgodb.Close()
	log.Release("closed")
}
