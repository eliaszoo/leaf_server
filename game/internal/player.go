package internal

import (
	"github.com/name5566/leaf/util"
	"leaf_server/base"
	"leaf_server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

const (
	INTATTR_COIN  		= 0 //金币
	INTATTR_POWER 		= 1 //体力

	INTATTR_MAX 		= 10
)

const (
	STRATTR_NICK = 0 //昵称
	STRATTR_ICON = 1 //头像

	STRATTR_MAX = 2
)


type Player struct {
	objid   		string
	agent   		gate.Agent
	UID 			int64				
	Account 		string						
	IntAttr 		[]int    //整型属性
	StrAttr 		[]string //字符串属性
}

func CreatePlayer() *Player {
	player := new(Player)
	player.IntAttr = make([]int, INTATTR_MAX)
	player.StrAttr = make([]string, STRATTR_MAX)
	return player
}

func (self *Player) GetIntAttr(index int) int {
	if index < 0 || index >= INTATTR_MAX {
		return 0
	}
	return self.IntAttr[index]
}

func (self *Player) SetIntAttr(index, val int) {
	if index < 0 || index >= INTATTR_MAX {
		return
	}

	self.IntAttr[index] = val
}

func (self *Player) GetStrAttr(index int) string {
	if index < 0 || index >= STRATTR_MAX {
		return ""
	}
	return self.StrAttr[index]
}

func (self *Player) SetStrAttr(index int, val string) {
	if index < 0 || index >= STRATTR_MAX {
		return
	}

	self.StrAttr[index] = val
}

func (self *Player) InitData(account string) {
	self.Account = account
	self.UID     = uidbuilder.GenerateUID()
}

func (self *Player) CallClientFunc(ret int, cmd string, ans interface{}) {
	errmsg := ""

	message := &msg.RetMsg{ret, errmsg, cmd, ans}
	self.agent.WriteMsg(message)
}

func (self *Player) Test(req *msg.TestReq) {

}

//登陆
func (self *Player) OnLogin() {
}

func (self *Player) OnLogout() {
	self.Save()

	TimerManager.RmvAllTimer(self)
}

//保存玩家数据
func (self *Player) Save() {
	mgodb.Set(base.DBTask{self.objid, base.DBNAME, base.PLAYERSET, "_id", base.BsonObjectID(self.objid), util.DeepClone(self), func(param interface{}, err error) {
		if nil != err {
			log.Error("save playerdata failed:", self.objid)
		}
	} })
}

//同步保存玩家数据
func (self *Player) SaveSync() {
	if nil != mgodb.SetSync(base.DBNAME, base.PLAYERSET, "_id", base.BsonObjectID(self.objid), self) {
		log.Error("save playerdata failed:", self.objid)
	}
}