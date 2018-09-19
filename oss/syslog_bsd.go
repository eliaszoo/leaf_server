// +build !windows

package oss

import (
	"log/syslog"
)

func New(flag Priority, tag string) (writer, error) {
	return syslog.New(syslog.Priority(flag), tag)
}

func Dial(net, addr string, flag Priority, tag string) (writer, error) {
	return syslog.Dial(net, addr, syslog.Priority(flag), tag)
}

