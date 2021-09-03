package ilua

import (
	"errors"
	"graduate/src/ilua/lmysql"
	"sync"
	lua "github.com/yuin/gopher-lua"
)


type LStatePool struct {
	mu sync.RWMutex
	ch chan *LState
	proto *lua.FunctionProto
}
type LState struct {
	name string
	*lua.LState
}

func (l * LState) CallByParam(cp lua.P, args ...lua.LValue) error {
	if l.LState == nil {
		return errors.New("lua state is nil")
	}
	err := l.LState.CallByParam(cp, args...)
	if err != nil {
		return err
	}
	return nil
}

//创建LState
func NewLState(proto *lua.FunctionProto) *LState {
	l := &LState{
		name: proto.SourceName,
	}
	l.LState = lua.NewState(
		lua.Options{
			RegistrySize: 1024 * 20,
			RegistryMaxSize: 1024 * 80,
			RegistryGrowStep: 120,
			CallStackSize: 120,
			IncludeGoStackTrace: true,
		},
		)
	l.RegisterModule("mysql", lmysql.Exports)

	lfunc := l.NewFunctionFromProto(proto)
	l.Push(lfunc)
	err := l.PCall(0, lua.MultRet, nil)
	if err != nil {
		l.Close()
		return nil
	}
	return l
}

const (
	LStatePoolSizeMax = 100
	LStatePoolSizeInit = 20
)

//创建一个lua VM 池
func NewLStatePool(proto *lua.FunctionProto) *LStatePool {
	p := &LStatePool{}
	p.ch = make(chan *LState, LStatePoolSizeMax)
	p.proto = proto
	for i := 0; i < LStatePoolSizeInit; i++ {
		l := NewLState(proto)
		if l == nil {
			return nil
		}
		p.ch <- l
	}
	return p
}

func (p *LStatePool) Get() *LState {
	p.mu.RLock()
	ch := p.ch
	p.mu.RUnlock()

	select {
	case l := <-ch:
		return l
	default:
		return NewLState(p.proto)
	}
}

func (p *LStatePool) Put(l *LState) {
	if l == nil {
		return
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.ch == nil {
		l.Close()
		return
	}
	select {
	case p.ch <- l:
		return
	default:
		l.Close()
	}
	return
}

func (p *LStatePool) Shutdown() {
	p.mu.RLock()
	ch := p.ch
	p.ch = nil
	p.mu.Unlock()
	close(ch)
	for l := range ch {
		l.Close()
	}
	return
}