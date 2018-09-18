package oss

import (
	"time"
	"fmt"
)

type Logger struct {
	writer
}

var gLogger *Logger

func NewLogger(network, addr, tag string) (*Logger, error) {
	var w writer
	var err error
	if "" == addr {
		w, err = New(LOG_INFO|LOG_USER, tag)
	} else {
		w, err = Dial(network, addr, LOG_INFO|LOG_USER, tag)
	}
	if nil != err {
		return nil, err
	}
	return &Logger{w}, err
}

func Export(logger *Logger) {
	if nil != logger {
		gLogger = logger
	}
}

func (logger *Logger) WriteLog(objid string, uid int64, billid int, param string) {
	now := time.Now()
	tm := fmt.Sprintf("%d%02d%02d_%02d%02d%02d",
	now.Year(),
	now.Month(),
	now.Day(),
	now.Hour(),
	now.Minute(),
	now.Second())

	gLogger.Info(fmt.Sprintf("%d^%s^%s^%d^%s", uid, tm, objid, billid, param))
}

func Close() {
	gLogger.Close()
}

//资源日志
func ResLog(objid string, uid int64, restyp, id, count, remaincount, reason int) {
	gLogger.WriteLog(objid, uid, BILLID_RES, fmt.Sprintf("%d|%d|%d|%d|%d", restyp, id, count, remaincount, reason))
}

//行为日志
func ActionLog(objid string, uid int64, actionid int, param interface{}) {
	str := ""
	if nil != param {
		str = fmt.Sprintf("%+v", param)
	}
	gLogger.WriteLog(objid, uid, actionid, str)
}