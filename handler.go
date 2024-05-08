package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var Handlers = map[string]func([]Value) Value{
	"PING":      ping,
	"SET":       set,
	"GET":       get,
	"DEL":       del,
	"RANDOMKEY": randkey,
	"EXISTS":    exists,
	"STRLEN":    strlen,
	"LOLWUT":    lolwut,
}

func ping(args []Value) Value {
	return Value{typ: "bulk", bulk: "PONG"}
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

	return Value{typ: "bulk", bulk: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "incorrect number of arg for GET command"}
	}
	var val string

	setMU.RLock()
	val, ok := setData[args[0].bulk]
	setMU.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
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

func randkey(args []Value) Value {
	if len(setData) == 0 {
		return Value{typ: "null"}
	}

	randNum := rand.Intn(len(setData))
	var randKey string
	counter := 0
	setMU.RLock()
	for key, _ := range setData {
		if counter == randNum {
			randKey = key
			break
		}
		counter++
	}
	setMU.RUnlock()
	return Value{typ: "bulk", bulk: randKey}
}

func exists(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "incorrect number of args for exisits command"}
	}
	counter := 0
	for _, a := range args {
		setMU.RLock()
		_, ok := setData[a.bulk]
		setMU.RUnlock()
		if ok {
			counter++
		}
	}
	return Value{typ: "integer", num: counter}
}

func strlen(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "incorrect number of args for strlen command"}
	}
	setMU.RLock()
	val := setData[args[0].bulk]
	setMU.RUnlock()

	return Value{typ: "integer", num: len(val)}

}

func lolwut(args []Value) Value {
	return Value{typ: "bulk", bulk: "GoKV 0.1 :):):)\r\n"}
}
