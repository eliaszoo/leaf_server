package internal

import (
	"leaf_server/base"
	"gopkg.in/mgo.v2"
	"github.com/name5566/leaf/log"
)

type UidBuilder struct {
	Seq 	int64
}

func (builder *UidBuilder) Init() {
	err := mgodb.GetSync(base.DBNAME, base.INCREMENTSET, "_id", "uid", builder)
	if mgo.ErrNotFound == err {
		mgodb.EnsureCounter(base.DBNAME, base.INCREMENTSET, "uid")
	} else if nil != err {
		log.Error("load uid seq failed:", err.Error())
		return
	}

	playercount := int64(mgodb.GetTableCountSync(base.DBNAME, base.PLAYERSET))
	if builder.Seq <= playercount {
		builder.Seq = playercount + 1
		mgodb.SetSync(base.DBNAME, base.INCREMENTSET, "_id", "uid", builder)
	}
}

func (builder *UidBuilder) GenerateUID() int64 {
	mgodb.IncreSeq(base.DBNAME, base.INCREMENTSET, "uid", func(i interface{}, err error) {
		if nil != err {
			log.Error("Increment uid failed:", err.Error())
		}
	})

	builder.Seq ++
	return 10000 + builder.Seq
}
