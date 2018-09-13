package base

import "gopkg.in/mgo.v2/bson"

const (
	DBNAME = "Test"

	ACCOUNTSET   = "AccountSet"  	//账号集合
	PLAYERSET  	 = "PlayerSet"	 	//玩家集合
	INCREMENTSET = "Increment"   	//增量集合
) 

type DBTask struct {
	ObjID 	   string
	DB         string
	Collection string
	Key        string
	KeyV       interface{}
	Ret        interface{}
	Cb 		   func(interface{}, error)
}

type DBSearch struct {
	ObjID 		string
	DB 			string
	Collection 	string
	M 			bson.M
	Limit 		int
	Skip 		int
	Ret 		interface{}
	Cb 			func(interface{}, error)
}

//账号信息
type AccountInfo struct {
	Account  string
	Password string
	ObjID    string
}

//转换一个Bson
func BsonObjectID(s string) bson.ObjectId {
	if s == "" {
		return bson.NewObjectId()
	}

	if bson.IsObjectIdHex(s) {
		return bson.ObjectIdHex(s)
	}

	return bson.ObjectId(s)
}
