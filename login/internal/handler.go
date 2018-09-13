package internal

import (
	"reflect"
	"leaf_server/base"
	"leaf_server/game"
	"leaf_server/msg"

	"github.com/goinggo/mapstructure"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"gopkg.in/mgo.v2/bson"
)

type clientFunc func(interface{}, gate.Agent) int

type clientFS struct {
	clientFunc
	req interface{}
}

var handleMap map[string]clientFS

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMap = make(map[string]clientFS)
	handleMsg(&msg.LoginMsg{}, onClientMsg)

	registerFunc("login", login, msg.LoginReq{})
}

func registerFunc(cmd string, fun clientFunc, req interface{}) {
	_, ok := handleMap[cmd]
	if ok {
		log.Error("login handle register repeated function:", cmd)
		return
	}

	handleMap[cmd] = clientFS{fun, req}
}

func done(agent gate.Agent, message *msg.RetMsg) {
	/*data, err := json.Marshal(message)
	if nil != err {
		log.Error("login module json marshal failed:", err.Error())
		return
	}*/

	agent.WriteMsg(message)
}

func onClientMsg(args []interface{}) {
	m := args[0].(*msg.LoginMsg)
	a := args[1].(gate.Agent)

	fs, ok := handleMap[m.Cmd]
	if !ok {
		log.Error("login function not found:", m.Cmd)
		return
	}

	i := reflect.New(reflect.TypeOf(fs.req)).Interface()
	err := mapstructure.Decode(m.Req, i)
	if nil != err {
		log.Error("login function decode req failed:", m.Cmd)
	}

	ret := fs.clientFunc(i, a)
	if 0 != ret {
		log.Error("login function execute error:", m.Cmd, " ", ret)
		return
	}

	/*
		message := new(msg.RetMsg)
		message.Ret = ret
		message.Cmd = m.Cmd
		message.Error = errmsg
		message.Ans = ans

		data, err := json.Marshal(message)
		if nil != err {
			log.Error("login module json marshal failed:", m.Cmd, " ", err.Error())
			return
		}
		a.WriteMsg(data)*/
}

func login(message interface{}, agent gate.Agent) int {
	req := message.(*msg.LoginReq)
	mgodb.Get(base.DBTask{req.Account, base.DBNAME, base.ACCOUNTSET, "account", req.Account, &base.AccountInfo{}, func(param interface{}, err error) {
		info := param.(*base.AccountInfo)
		if info.Account == "" {
			info.Account = req.Account
			info.Password = req.Password
			info.ObjID = bson.NewObjectId().Hex()
			mgodb.Set(base.DBTask{info.Account, base.DBNAME, base.ACCOUNTSET, "account", req.Account, info, nil})
		}

		if info.Password != req.Password {
			done(agent, &msg.RetMsg{1, "", "login", nil})
			return
		}

		agent.SetUserData(info)
		skeleton.AsynCall(game.ChanRPC, "LoginSuccess", agent, func(err error) {
			if nil != err {
				log.Error("login failed:", info.ObjID, " ", err.Error())
				done(agent, &msg.RetMsg{-1, "", "login", nil})
				return
			}
			//done(agent, &msg.RetMsg{0, "", "login", &msg.LoginAns{info.ObjID}})
		})		
	} })
	return 0
}
