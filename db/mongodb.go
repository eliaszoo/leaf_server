package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/name5566/leaf/go"
	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"

	"leaf_server/base"
)

type Mongodb struct {
	*mongodb.DialContext
	linearWorkers 	[]*g.LinearContext
}

func Dial(url string, sessionNum int, skeleton *module.Skeleton) (*Mongodb, error) {
	c, err := mongodb.Dial(url, sessionNum)
	if nil != err {
		log.Error("dial mongodb failed:", url, " ", err.Error())
		return nil, err
	}
	if nil == skeleton {
		return &Mongodb{c, nil}, err
	}

	workers := make([]*g.LinearContext, 0)
	for i := 0; i < sessionNum>>1; i ++ {
		workers = append(workers, skeleton.NewLinearContext())
	}
	return &Mongodb{c, workers}, err
}

func hash(str string) int {
	h := uint(1315423911)
	for i := 0; i < len(str); i ++ {
		h ^= ((uint(str[i]) << 5) + uint(str[i]) + (h >> 2))
	}
	return int(h & 0x7fffffff)
}

func (db *Mongodb) worker(objid string) *g.LinearContext {
	index := hash(objid) % len(db.linearWorkers)
	return db.linearWorkers[index]
}

func (db *Mongodb) Get(task base.DBTask) {
	w := db.worker(task.ObjID)
	var err error
	w.Go(func() {
		s := db.Ref()
		defer db.UnRef(s)
		
		err = s.DB(task.DB).C(task.Collection).Find(bson.M{task.Key: task.KeyV}).One(task.Ret)
		if nil != err && mgo.ErrNotFound != err {
			log.Error("db get failed:", task.DB, " collection:", task.Collection, " key:", task.Key, " error:", err.Error())
		}
	}, func() {
		if nil != task.Cb {
			task.Cb(task.Ret, err)
		}
	})
}

func (db *Mongodb) GetAll(task base.DBTask) {
	w := db.worker(task.ObjID)
	var err error
	w.Go(func() {
		s := db.Ref()
		defer db.UnRef(s)
		
		err = s.DB(task.DB).C(task.Collection).Find(bson.M{task.Key: task.KeyV}).All(task.Ret)
		if nil != err && mgo.ErrNotFound != err {
			log.Error("db get failed:", task.DB, " collection:", task.Collection, " key:", task.Key, " error:", err.Error())
		} 
	}, func() {
		if nil != task.Cb {
			task.Cb(task.Ret, err)
		}
	})
}

func (db *Mongodb) Set(task base.DBTask) {
	w := db.worker(task.ObjID)
	var err error
	w.Go(func() {
		s := db.Ref()
		defer db.UnRef(s)
		
		_, err = s.DB(task.DB).C(task.Collection).Upsert(bson.M{task.Key: task.KeyV}, task.Ret)
		if nil != err {
			log.Error("db set failed:", task.DB, " collection:", task.Collection, " key:", task.Key, " error:", err.Error())
		}
	}, func() {
		if nil != task.Cb {
			task.Cb(task.Ret, err)
		}
	})
}

func (db *Mongodb) IncreSeq(dbname, collection, id string, cb func(interface{}, error)) {
	w := db.worker(id)
	var res struct { Seq int64 }
	var err error
	w.Go(func() {
		s := db.Ref()
		defer db.UnRef(s)
		
		_, err = s.DB(dbname).C(collection).FindId(id).Apply(mgo.Change{
			Update:    bson.M{"$inc": bson.M{"seq": 1}},
			ReturnNew: true,
		}, &res)
	}, func() {
		if nil != cb {
			cb(res.Seq, err)
		}
	})
}

func (db *Mongodb) GetTableCount(objid, dbname, collection string, cb func(interface{}, error)) {
	w := db.worker(objid)
	count := 0
	var err error
	w.Go(func() {
		s := db.Ref()
		defer db.UnRef(s)
		
		count, err = s.DB(dbname).C(collection).Count()
	}, func() {
		if nil != cb {
			cb(count, err)
		}
	})
}

func (db *Mongodb) Search(task base.DBSearch) {
	w := db.worker(task.ObjID)
	var err error
	w.Go(func() {
		s := db.Ref()
		defer db.UnRef(s)
		
		err = s.DB(task.DB).C(task.Collection).Find(task.M).Limit(task.Limit).Skip(task.Skip).All(task.Ret)
		if nil != err && mgo.ErrNotFound != err {
			log.Error("db search failed:", task.DB, " collection:", task.Collection, " error:", err.Error())
		}
	}, func() {
		if nil != task.Cb {
			task.Cb(task.Ret, err)
		}
	})
}

func (db *Mongodb) GetSync(dbname, collection, key string, keyv, ret interface{}) error {
	s := db.Ref()
	defer db.UnRef(s)
	
	err := s.DB(dbname).C(collection).Find(bson.M{key: keyv}).One(ret)
	if nil != err && mgo.ErrNotFound != err {
		log.Error("db get failed:", dbname, " collection:", collection, " key:", key, " error:", err.Error())
	}
	return err
}

func (db *Mongodb) SetSync(dbname, collection, key string, keyv, v interface{}) error {
	s := db.Ref()
	defer db.UnRef(s)
	
	_, err := s.DB(dbname).C(collection).Upsert(bson.M{key: keyv}, v)
	if nil != err {
		log.Error("db set failed:", dbname, " collection:", collection, " key:", key, " error:", err.Error())
	}
	return err
}

func (db *Mongodb) GetTableCountSync(dbname, collection string) int {
	s := db.Ref()
	defer db.UnRef(s)
	count, err := s.DB(dbname).C(collection).Count()
	if err != nil {
		return 0
	}
	return count
}

func (db *Mongodb) SearchSync(dbname, collection string, m bson.M, ret interface{}, limit, skip int) error {
	s := db.Ref()
	defer db.UnRef(s)
	
	err := s.DB(dbname).C(collection).Find(m).Limit(limit).Skip(skip).All(ret)
	if nil != err && mgo.ErrNotFound != err {
		log.Error("db search failed:", dbname, " collection:", collection, " error:", err.Error())
	}
	return err
}