package activity

import (
	"container/list"
	"crypto/md5"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"graduate/src/event"
	"graduate/src/ilua"
	"graduate/src/session"
	luar "layeh.com/gopher-luar"
	"sync"
)

var actCenter sync.Map

type fileActivity struct {
	Act *Activity
	File string
}

type Activity struct {
	ID string
	Rule string
	RuleMD5 [16]byte
	EventsList list.List
	LuaVMPool *ilua.LStatePool
}

func Get(id string) *Activity {
	v, ok := actCenter.Load(id)
	if !ok {
		return nil
	}
	if v == nil {
		return nil
	}
	return v.(fileActivity).Act
}

func Del(id string) {
	act := Get(id)
	if act == nil{
		return
	}
	actCenter.Delete(id)
	act.DeInit()
	return
}

func Update(id string, rule string, name string) (*Activity, error) {
	newAct := new(id, rule)
	af := fileActivity{
		Act: newAct,
		File: name,
	}
	oldActExist := false
	if act := Get(id); act != nil {
		if act.RuleMD5 == newAct.RuleMD5 {
			return act, nil
		}
		oldActExist = true
	}
	err := newAct.init()
	if err != nil {
		return nil, err
	}
	if oldActExist {
		Del(id)
	}
	actCenter.Store(id, af)
	return newAct, nil
}

func (a *Activity) Update(eventType event.Type, eventParam event.Param) {
	l := a.LuaVMPool.Get()
	defer a.LuaVMPool.Put(l)
	if nil == l {
		return
	}
	err := l.CallByParam(
		lua.P{
			Fn: l.GetGlobal("Update"),
			NRet: 0,
			Protect: true,
			Handler: nil,
		},
		lua.LNumber(eventType),
		luar.New(l.LState, map[string]interface{}(eventParam)),
		)
	if err != nil {
		return
	}
	return
}

func (a *Activity) QueryDetails(s *session.Session, req []byte) ([]byte, error) {
	l := a.LuaVMPool.Get()
	defer a.LuaVMPool.Put(l)
	err := l.CallByParam(
		lua.P{
			Fn: l.GetGlobal("QueryDetail"),
			NRet: 0,
			Protect: true,
			Handler: nil,
		},
		luar.New(l.LState, s),
		lua.LString(req),
	)
	if err != nil {
		return nil, err
	}
	lres := l.Get(-1)
	l.Pop(1)
	res := lua.LVAsString(lres)
	return []byte(res), nil
}

func (a *Activity) remove(t event.Type) {
	event.Event[t].Detach(a)
}

func new(id string, rule string) *Activity {
	act := &Activity{}
	act.Rule = rule
	act.ID = id
	act.RuleMD5 = md5.Sum([]byte(rule))
	return act
}

func (a *Activity) init() error {
	luaName := fmt.Sprintf("act_%v_%x.lua", a.ID, a.RuleMD5[0:3])
	proto, _ := ilua.Compile(luaName, a.Rule)
	a.LuaVMPool = ilua.NewLStatePool(proto)
	l := a.LuaVMPool.Get()
	defer a.LuaVMPool.Put(l)

	events := l.GetGlobal("Events")
	key, value := events.(*lua.LTable).Next(lua.LNil)
	for key != lua.LNil {
		v := lua.LVAsNumber(value)
		t := event.Type(v)
		a.register(t)
		a.EventsList.PushBack(v)
		key, value = events.(*lua.LTable).Next(key)
	}
	return nil
}

func (a *Activity) register(t event.Type) {
	event.Event[t].Attach(a)
}


func (a *Activity) DeInit() {
	for e := a.EventsList.Front(); e != nil; e = e.Next() {
		t := e.Value.(event.Type)
		a.remove(t)
	}
	a.LuaVMPool.Shutdown()
	return
}

