package internal

import (
	"leaf_server/base"
	"leaf_server/db"
	"leaf_server/conf"

	"github.com/name5566/leaf/module"
	"github.com/name5566/leaf/log"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	mgodb    = new(db.Mongodb)
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	var err error
	mgodb, err = db.Dial(conf.Server.MgodbAddr, conf.Server.LoginMgoConnNum, skeleton)
	if nil == mgodb {
		log.Error("dial mongodb failed:", conf.Server.MgodbAddr, " ", err.Error())
	}

	mgodb.EnsureUniqueIndex(base.DBNAME, base.ACCOUNTSET, []string{"account"})
}

func (m *Module) OnDestroy() {
	mgodb.Close()
}
