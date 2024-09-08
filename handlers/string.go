package handlers

import (
	"strconv"

	"math/rand"

	"github.com/idugan100/GoKV/resp"
)

func set(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "SET"}.Error()}
	}
	key := args[0].Bulk
	val := args[1].Bulk
	stringMU.Lock()
	stringData[key] = stringDataItem{str: val, will_expire: false}
	stringMU.Unlock()

	return resp.Serializable{Typ: "bulk", Bulk: "OK"}
}

func get(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "GET"}.Error()}
	}

	str, ok := getString(args[0].Bulk)

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	return resp.Serializable{Typ: "bulk", Bulk: str}
}

func del(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "DEL"}.Error()}
	}
	deletedCounter := 0
	for i := 0; i < len(args); i++ {
		_, ok := getString(args[i].Bulk)
		if ok {
			deletedCounter++
			stringMU.Lock()
			delete(stringData, args[i].Bulk)
			stringMU.Unlock()
		}
	}
	return resp.Serializable{Typ: "integer", Num: deletedCounter}
}

func randkey(args []resp.Serializable) resp.Serializable {
	if len(stringData) == 0 {
		return resp.Serializable{Typ: "null"}
	}

	randNum := rand.Intn(len(stringData))
	var randKey string
	counter := 0
	stringMU.RLock()
	for key := range stringData {
		if counter == randNum {
			randKey = key
			break
		}
		counter++
	}
	stringMU.RUnlock()
	return resp.Serializable{Typ: "bulk", Bulk: randKey}
}

func exists(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "EXISTS"}.Error()}
	}
	counter := 0

	stringMU.RLock()
	defer stringMU.RUnlock()

	for _, a := range args {
		_, ok := getString(a.Bulk)
		if ok {
			counter++
		}
	}
	return resp.Serializable{Typ: "integer", Num: counter}
}

func strlen(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "STRLEN"}.Error()}
	}
	stringMU.RLock()
	val, ok := getString(args[0].Bulk)
	stringMU.RUnlock()
	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	return resp.Serializable{Typ: "integer", Num: len(val)}

}

func getset(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "GETSET"}.Error()}
	}
	oldSerializable, ok := getString(args[0].Bulk)
	stringMU.Lock()
	stringData[args[0].Bulk] = stringDataItem{str: args[1].Bulk, will_expire: false}
	stringMU.Unlock()

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	return resp.Serializable{Typ: "bulk", Bulk: oldSerializable}

}

func setnx(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "SETNX"}.Error()}
	}

	_, ok := getString(args[0].Bulk)

	if ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	stringMU.Lock()
	defer stringMU.Unlock()
	stringData[args[0].Bulk] = stringDataItem{str: args[1].Bulk, will_expire: false}
	return resp.Serializable{Typ: "integer", Num: 1}

}

func incr(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "INCR"}.Error()}
	}

	str, ok := getString(args[0].Bulk)

	if !ok {
		stringMU.Lock()
		stringData[args[0].Bulk] = stringDataItem{str: "0", will_expire: false}
		stringMU.Unlock()
		str = "0"
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "INCR"}.Error()}
	}
	num++

	stringMU.Lock()
	stringData[args[0].Bulk] = stringDataItem{str: strconv.Itoa(int(num)), will_expire: false}
	stringMU.Unlock()

	return resp.Serializable{Typ: "integer", Num: int(num)}
}

func incrby(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "INCRBY"}.Error()}
	}

	str, ok := getString(args[0].Bulk)

	if !ok {
		stringMU.Lock()
		stringData[args[0].Bulk] = stringDataItem{str: "0", will_expire: false}
		stringMU.Unlock()
		str = "0"
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "INCRBY"}.Error()}
	}

	incrementAmount, err := strconv.ParseInt(args[1].Bulk, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "INCRBY"}.Error()}
	}
	num += incrementAmount
	stringMU.Lock()
	stringData[args[0].Bulk] = stringDataItem{str: strconv.Itoa(int(num)), will_expire: false}
	stringMU.Unlock()

	return resp.Serializable{Typ: "integer", Num: int(num)}
}

func decr(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "DECR"}.Error()}
	}

	str, ok := getString(args[0].Bulk)

	if !ok {
		stringMU.Lock()
		stringData[args[0].Bulk] = stringDataItem{str: "0", will_expire: false}
		stringMU.Unlock()
		str = "0"
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "DECR"}.Error()}
	}
	num--
	stringMU.Lock()
	stringData[args[0].Bulk] = stringDataItem{str: strconv.Itoa(int(num)), will_expire: false}
	stringMU.Unlock()
	return resp.Serializable{Typ: "integer", Num: int(num)}
}

func decrby(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "DECRBY"}.Error()}
	}

	str, ok := getString(args[0].Bulk)

	if !ok {
		stringMU.Lock()
		stringData[args[0].Bulk] = stringDataItem{str: "0", will_expire: false}
		stringMU.Unlock()
		str = "0"
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "DECRBY"}.Error()}
	}

	decrementAmount, err := strconv.ParseInt(args[1].Bulk, 10, 64)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "DECRBY"}.Error()}
	}
	num -= decrementAmount
	stringMU.Lock()
	stringData[args[0].Bulk] = stringDataItem{str: strconv.Itoa(int(num)), will_expire: false}
	stringMU.Unlock()

	return resp.Serializable{Typ: "integer", Num: int(num)}
}

func renamenx(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "RENAMENX"}.Error()}
	}

	_, ok := getString(args[0].Bulk)
	if !ok {
		return resp.Serializable{Typ: "error", Str: "key to be renamed not found"}
	}

	//return 0 if new key already exisits
	_, ok = getString(args[1].Bulk)
	if ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}

	stringMU.Lock()
	stringData[args[1].Bulk] = stringData[args[0].Bulk]
	delete(stringData, args[0].Bulk)
	stringMU.Unlock()
	return resp.Serializable{Typ: "integer", Num: 1}
}

func rename(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "RENAME"}.Error()}
	}

	_, ok := getString(args[0].Bulk)
	if !ok {
		return resp.Serializable{Typ: "error", Str: "key to be renamed not found"}
	}
	stringMU.Lock()
	stringData[args[1].Bulk] = stringData[args[0].Bulk]
	delete(stringData, args[0].Bulk)
	stringMU.Unlock()
	return resp.Serializable{Typ: "bulk", Bulk: "OK"}
}

func mget(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "MGET"}.Error()}
	}

	results := []resp.Serializable{}
	for _, arg := range args {
		val, ok := getString(arg.Bulk)
		if !ok {
			results = append(results, resp.Serializable{Typ: "null"})
			continue
		} else {
			results = append(results, resp.Serializable{Typ: "bulk", Bulk: val})
		}
	}
	return resp.Serializable{Typ: "array", Array: results}
}
