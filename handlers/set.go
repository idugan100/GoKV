package handlers

import (
	"math/rand"
	"strconv"

	"github.com/idugan100/GoKV/resp"
)

func set(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arg for SET command"}
	}
	key := args[0].Bulk
	val := args[1].Bulk
	setMU.Lock()
	setData[key] = val
	setMU.Unlock()

	return resp.Serializable{Typ: "bulk", Bulk: "OK"}
}

func get(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arg for GET command"}
	}
	var val string

	setMU.RLock()
	val, ok := setData[args[0].Bulk]
	setMU.RUnlock()
	if !ok {
		return resp.Serializable{Typ: "null"}
	}
	return resp.Serializable{Typ: "bulk", Bulk: val}
}

func del(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of args for DEL command"}
	}
	deletedCounter := 0
	setMU.Lock()
	for i := 0; i < len(args); i++ {
		_, ok := setData[args[i].Bulk]
		if ok {
			deletedCounter++
			delete(setData, args[i].Bulk)
		}
	}
	setMU.Unlock()
	return resp.Serializable{Typ: "integer", Num: deletedCounter}
}

func randkey(args []resp.Serializable) resp.Serializable {
	if len(setData) == 0 {
		return resp.Serializable{Typ: "null"}
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
	return resp.Serializable{Typ: "bulk", Bulk: randKey}
}

func exists(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of args for exisits command"}
	}
	counter := 0
	for _, a := range args {
		setMU.RLock()
		_, ok := setData[a.Bulk]
		setMU.RUnlock()
		if ok {
			counter++
		}
	}
	return resp.Serializable{Typ: "integer", Num: counter}
}

func strlen(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of args for Strlen command"}
	}
	setMU.RLock()
	val := setData[args[0].Bulk]
	setMU.RUnlock()

	return resp.Serializable{Typ: "integer", Num: len(val)}

}

func getset(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements for GETSET command"}
	}
	setMU.Lock()
	oldSerializable, ok := setData[args[0].Bulk]
	setData[args[0].Bulk] = args[1].Bulk
	setMU.Unlock()

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	return resp.Serializable{Typ: "bulk", Bulk: oldSerializable}

}

func setnx(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements for SETNX command"}
	}
	setMU.RLock()
	_, ok := setData[args[0].Bulk]
	setMU.RUnlock()

	if ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}

	setMU.Lock()
	setData[args[0].Bulk] = args[1].Bulk
	setMU.Unlock()
	return resp.Serializable{Typ: "integer", Num: 1}

}

func incr(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements for INCR command"}
	}
	setMU.Lock()
	defer setMU.Unlock()

	val, ok := setData[args[0].Bulk]

	if !ok {
		setData[args[0].Bulk] = "0"
		val = "0"
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: "incorrect data type for INCR operation"}
	}
	num++
	setData[args[0].Bulk] = strconv.Itoa(int(num))

	return resp.Serializable{Typ: "integer", Num: int(num)}
}
func decr(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements for DECR command"}
	}
	setMU.Lock()
	defer setMU.Unlock()

	val, ok := setData[args[0].Bulk]

	if !ok {
		setData[args[0].Bulk] = "0"
		val = "0"
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: "incorrect data type for DECR operation"}
	}
	num--
	setData[args[0].Bulk] = strconv.Itoa(int(num))

	return resp.Serializable{Typ: "integer", Num: int(num)}
}

func decrby(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguments for DECRBY command"}
	}
	setMU.Lock()
	defer setMU.Unlock()

	val, ok := setData[args[0].Bulk]

	if !ok {
		setData[args[0].Bulk] = "0"
		val = "0"
	}

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: "incorrect data type for DECRBY operation"}
	}

	decrementAmount, err := strconv.ParseInt(args[1].Bulk, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: "incorrect data type for decrement amount"}
	}
	num -= decrementAmount
	setData[args[0].Bulk] = strconv.Itoa(int(num))

	return resp.Serializable{Typ: "integer", Num: int(num)}
}

func rename(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements for RENAME command"}
	}
	setMU.Lock()
	defer setMU.Unlock()

	val, ok := setData[args[0].Bulk]
	if !ok {
		return resp.Serializable{Typ: "error", Str: "key to be renamed not found"}
	}

	delete(setData, args[0].Bulk)
	setData[args[1].Bulk] = val
	return resp.Serializable{Typ: "bulk", Bulk: "OK"}
}
