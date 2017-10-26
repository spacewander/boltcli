package main

import (
	"strings"

	"github.com/Shopify/go-lua"
)

var (
	vm *lua.State

	curCmd string
)

func init() {
	vm = lua.NewState()
	lua.OpenLibraries(vm)
	injectAPI(vm)
}

func injectAPI(L *lua.State) {
	L.CreateTable(0, 1)

	L.CreateTable(0, 1)
	L.PushGoFunction(dispatchCmd)
	L.SetField(-2, "__index")
	L.SetMetaTable(-2)

	// inject global api namespace
	L.Global("package")
	L.Field(-1, "loaded")
	L.PushValue(-3)
	L.SetField(-2, "bolt")
	L.Pop(2)
	L.SetGlobal("bolt")
}

func dispatchCmd(L *lua.State) int {
	// ignore the meta table itself (the first arg)
	if s, ok := lua.ToStringMeta(L, 2); ok {
		s = strings.ToLower(s)
		_, ok = CmdMap[s]
		if ok {
			curCmd = s
			L.PushGoFunction(execCmdInLuaScript)
			return 1
		}
	}
	// it is equal to return nil
	return 0
}

func pushList(L *lua.State, res []string) {
	L.CreateTable(len(res), 0)
	for i, s := range res {
		L.PushString(s)
		L.RawSetInt(-2, i+1)
	}
}

func pushMap(L *lua.State, res map[string]interface{}) {
	L.CreateTable(0, len(res))
	for k, v := range res {
		L.PushString(k)
		switch v := v.(type) {
		case int64:
			L.PushInteger(int(v))
		case string:
			L.PushString(v)
		case map[string]interface{}:
			pushMap(L, v)
		}
		L.RawSet(-3)
	}
}

func execCmdInLuaScript(L *lua.State) int {
	args := []string{}
	nargs := L.Top()
	for i := 1; i <= nargs; i++ {
		luaType := L.TypeOf(i)
		switch luaType {
		case lua.TypeNumber:
			fallthrough
		case lua.TypeString:
			if s, ok := lua.ToStringMeta(L, i); ok {
				args = append(args, s)
			}
		default:
			// arg x is one based, like other stuff in lua land
			L.PushFString("The type of arg %d is incorrect, only number and string are acceptable", i)
			L.Error()
		}
	}
	// we have checked the existence of 'curCmd' before
	f, _ := CmdMap[curCmd]
	res, err := f(args...)
	if err != nil {
		L.PushNil()
		L.PushString(err.Error())
		return 2
	}
	switch res := res.(type) {
	case bool:
		L.PushBoolean(res)
	case []byte:
		L.PushString(string(res))
	case string:
		L.PushString(res)
	case []string:
		pushList(L, res)
	case map[string]interface{}:
		pushMap(L, res)
	case int:
		L.PushInteger(res)
	default:
		L.PushFString("The type of result returns from command '%s' with args %v is unsupported", curCmd, args)
		L.Error()
	}
	return 1
}

// StartScript evals given script file
func StartScript(script string) error {
	return lua.DoFile(vm, script)
}
