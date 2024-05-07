package main

import (
	"fmt"
	"sync"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"DEL":  del,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "incorrect number of arg for SET command"}
	}
	key := args[0].bulk
	val := args[1].bulk
	setMU.Lock()
	setData[key] = val
	setMU.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "incorrect number of arg for GET command"}
	}
	var val string

	setMU.RLock()
	val = setData[args[0].bulk]
	setMU.RUnlock()
	return Value{typ: "string", str: val}
}

func del(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "incorrect number of args for DEL command"}
	}
	deletedCounter := 0
	setMU.Lock()
	for i := 0; i < len(args); i++ {
		_, ok := setData[args[i].bulk]
		if ok {
			deletedCounter++
			delete(setData, args[i].bulk)
		}
	}
	setMU.Unlock()
	fmt.Println("reached here deleted ", deletedCounter)
	return Value{typ: "integer", num: deletedCounter}
}
