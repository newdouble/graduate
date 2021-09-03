package ilua

import (
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"strings"
)

func Compile(name string, code string) (*lua.FunctionProto, error) {
	chunk, err := parse.Parse(strings.NewReader(code), name)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, name)
	return proto, nil
}
