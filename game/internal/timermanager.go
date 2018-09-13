package internal

import (
	"time"
	"github.com/name5566/leaf/log"
)

type Timer struct {
	key 		string
	cb 			func()
	interval 	int
	loop 		bool
}

type timerManager struct {
	timerMap map[interface{}][]*Timer
}

func NewTimerManager() *timerManager {
	return &timerManager{make(map[interface{}][]*Timer)}
}

func (m *timerManager) addTimer(obj interface{}, key string, cb func(), interval int, loop bool) {
	skeleton.AfterFunc(time.Millisecond * time.Duration(interval), func() {
		if !m.timerVaild(obj, key) {
			return
		}
		if !loop {
			m.RmvTimer(obj, key)
		}
		
		cb()
		if loop {
			m.addTimer(obj, key, cb, interval, loop)
		}
	})
}

func (m *timerManager) getTimerIndex(l []*Timer, key string) int {
	for i, timer := range l {
		if timer.key == key {
			return i
		}
	}
	return -1
}

func (m *timerManager) timerVaild(obj interface{}, key string) bool {
	l, ok := m.timerMap[obj]
	if !ok {
		return false
	}

	return m.getTimerIndex(l, key) >= 0
}

func (m *timerManager) AddTimer(obj interface{}, key string, interval int, loop bool, cb func()) {
	l, _ := m.timerMap[obj]
	if m.getTimerIndex(l, key) >= 0 {
		log.Error("add repeated timer:", key)
		return
	}

	l = append(l, &Timer{key, cb, interval, loop})
	m.timerMap[obj] = l
	m.addTimer(obj, key, cb, interval, loop)
}

func (m *timerManager) RmvTimer(obj interface{}, key string) {
	l, ok := m.timerMap[obj]
	if !ok {
		return 
	}

	index := m.getTimerIndex(l, key)
	if index >= 0 {
		l = append(l[:index], l[index+1:]...)
	}
	m.timerMap[obj] = l
}

func (m *timerManager) RmvAllTimer(obj interface{}) {
	delete(m.timerMap, obj)
}