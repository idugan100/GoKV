package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var Handlers = map[string]func([]Serializable) Serializable{
	"PING":      ping,
	"SET":       set,
	"GET":       get,
	"DEL":       del,
	"RANDOMKEY": randkey,
	"EXISTS":    exists,
	"STRLEN":    strlen,
	"LOLWUT":    lolwut,
	"FLUSHALL":  flushall,
	"GETSET":    getset,
}

func ping(args []Serializable) Serializable {
	return Serializable{typ: "bulk", bulk: "PONG"}
}

func set(args []Serializable) Serializable {
	if len(args) != 2 {
		return Serializable{typ: "error", str: "incorrect number of arg for SET command"}
	}
	key := args[0].bulk
	val := args[1].bulk
	setMU.Lock()
	setData[key] = val
	setMU.Unlock()

	return Serializable{typ: "bulk", bulk: "OK"}
}

func get(args []Serializable) Serializable {
	if len(args) != 1 {
		return Serializable{typ: "error", str: "incorrect number of arg for GET command"}
	}
	var val string

	setMU.RLock()
	val, ok := setData[args[0].bulk]
	setMU.RUnlock()
	if !ok {
		return Serializable{typ: "null"}
	}
	return Serializable{typ: "bulk", bulk: val}
}

func del(args []Serializable) Serializable {
	if len(args) < 1 {
		return Serializable{typ: "error", str: "incorrect number of args for DEL command"}
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
	return Serializable{typ: "integer", num: deletedCounter}
}

func randkey(args []Serializable) Serializable {
	if len(setData) == 0 {
		return Serializable{typ: "null"}
	}

	randNum := rand.Intn(len(setData))
	var randKey string
	counter := 0
	setMU.RLock()
	for key := range setData {
		if counter == randNum {
			randKey = key
			break
		}
		counter++
	}
	setMU.RUnlock()
	return Serializable{typ: "bulk", bulk: randKey}
}

func exists(args []Serializable) Serializable {
	if len(args) < 1 {
		return Serializable{typ: "error", str: "incorrect number of args for exisits command"}
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
	return Serializable{typ: "integer", num: counter}
}

func strlen(args []Serializable) Serializable {
	if len(args) != 1 {
		return Serializable{typ: "error", str: "incorrect number of args for strlen command"}
	}
	setMU.RLock()
	val := setData[args[0].bulk]
	setMU.RUnlock()

	return Serializable{typ: "integer", num: len(val)}

}

func lolwut(args []Serializable) Serializable {
	return Serializable{typ: "bulk", bulk: "GoKV 0.1 :):):)\r\n"}
}

func flushall(args []Serializable) Serializable {
	setMU.Lock()
	setData = map[string]string{}
	setMU.Unlock()
	return Serializable{typ: "string", str: "OK"}
}

func getset(args []Serializable) Serializable {
	if len(args) != 2 {
		return Serializable{typ: "error", str: "incorrect number of arguements for GETSET command"}
	}
	setMU.Lock()
	oldSerializable, ok := setData[args[0].bulk]
	setData[args[0].bulk] = args[1].bulk
	setMU.Unlock()

	if !ok {
		return Serializable{typ: "null"}
	}

	return Serializable{typ: "bulk", bulk: oldSerializable}

}
