package internal

import (
	"leaf_server/oss"
	"leaf_server/base"
	"leaf_server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("LoginSuccess", rpcLoginSuccess)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("new agent")
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("close agent")

	userdata := a.UserData()
	if nil == userdata {
		return
	}

	info := userdata.(*base.AccountInfo)
	player := PlayerManager.Get(info.ObjID)
	if nil != player {
		player.OnLogout()
	}
	PlayerManager.DelPlayer(info.ObjID)
}

func rpcLoginSuccess(args []interface{}) {
	a := args[0].(gate.Agent)
	userdata := a.UserData()
	info := userdata.(*base.AccountInfo)

	//判断玩家重复登陆
	player := PlayerManager.Get(info.ObjID)
	if nil != player {
		player.agent.Close()
		player.agent = a
		return
	}

	mgodb.Get(base.DBTask{info.ObjID, base.DBNAME, base.PLAYERSET, "_id", base.BsonObjectID(info.ObjID), CreatePlayer(), func(param interface{}, err error) {		
		player := param.(*Player)
		player.objid = info.ObjID
		player.agent = a
		if player.Account == "" { //保存新玩家数据
			player.InitData(info.Account)
			player.Save()
			oss.ActionLog(player.objid, player.UID, oss.BILLID_REGISTER, nil)
		}

		player.OnLogin()
		PlayerManager.AddPlayer(player)

		player.CallClientFunc(0, "login", &msg.LoginAns{info.ObjID})
		oss.ActionLog(player.objid, player.UID, oss.BILLID_LOGIN, nil)
	} })
}
