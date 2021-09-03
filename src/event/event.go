package event

import (
	"container/list"
	"sync"
)

type Subject interface {
	Attach(Observer)
	Detach(Observer)
	Notify(Type, Param)
}

type Observer interface {
	Update(Type, Param)
}


//Param 事件通知参数
type Param map[string]interface{}

type Type uint

func (p Param) Add(key string, value interface{}) {
	if p == nil {
		return
	}
	p[key] = value
	return
}
func (t Type) IsValid() bool {
	return t < EventMaxNumber
}

const (
	EventNewSelf Type = iota
	EventMaxNumber //事件总数

)

//事件主题对象映射表
var Event [EventMaxNumber]baseSubject

type baseSubject struct {
	rwm sync.RWMutex
	observers list.List
}

type baseObserver struct {
	subjects list.List
}

func (s *baseSubject) Attach(o Observer)  {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	for ob := s.observers.Front(); ob != nil; ob = ob.Next() {
		if ob.Value.(Observer) == o {
			return
		}
		s.observers.PushBack(o)
	}
}

func (s *baseSubject) Detach(o Observer) {
	s.rwm.Lock()
	defer s.rwm.Unlock()
	for ob := s.observers.Front(); ob != nil; ob = ob.Next() {
		if ob.Value.(Observer) == o {
			s.observers.Remove(ob)
		}
		break
	}
	return
}

func (s *baseSubject) Notify(t Type, p Param) {
	s.rwm.RLock()
	defer s.rwm.RUnlock()
	for ob := s.observers.Front(); ob != nil; ob = ob.Next() {
		ob.Value.(Observer).Update(t, p)
	}
}