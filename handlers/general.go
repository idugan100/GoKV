package handlers

import (
	"container/list"
	"fmt"
	"strconv"
	"time"

	"github.com/idugan100/GoKV/resp"
)

func ping(args []resp.Serializable) resp.Serializable {
	return resp.Serializable{Typ: "bulk", Bulk: "PONG"}
}

func lolwut(args []resp.Serializable) resp.Serializable {
	return resp.Serializable{Typ: "bulk", Bulk: "GoKV 0.1 :):):)\r\n"}
}

func flushall(args []resp.Serializable) resp.Serializable {
	stringMU.Lock()
	stringData = map[string]stringDataItem{}
	stringMU.Unlock()
	hstringMU.Lock()
	hsetData = map[string]map[string]string{}
	hstringMU.Unlock()
	listMU.Lock()
	listData = map[string]*list.List{}
	listMU.Unlock()
	return resp.Serializable{Typ: "string", Str: "OK"}
}

func dbsize(args []resp.Serializable) resp.Serializable {
	stringMU.RLock()
	defer stringMU.RUnlock()
	listMU.RLock()
	defer listMU.RUnlock()
	hstringMU.RLock()
	defer hstringMU.RUnlock()
	total_keys := len(stringData) + len(listData) + len(hsetData)

	return resp.Serializable{Typ: "integer", Num: total_keys}
}

func ttl(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "TTL"}.Error()}
	}
	fmt.Println("here 1")
	_, ok := getString(args[0].Bulk)
	fmt.Println("here 2")

	if !ok {
		return resp.Serializable{Typ: "integer", Num: -2}
	}
	fmt.Println("here 3")

	stringMU.RLock()
	defer stringMU.RUnlock()
	fmt.Println("here 4")

	val, ok := stringData[args[0].Bulk]
	fmt.Println("here 5")

	if !ok || !val.will_expire {
		return resp.Serializable{Typ: "integer", Num: -1}
	}
	fmt.Println("here 6")

	return resp.Serializable{Typ: "integer", Num: int(time.Until(val.expiration).Seconds())}
}

func pttl(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "PTTL"}.Error()}
	}
	_, ok := getString(args[0].Bulk)
	if !ok {
		return resp.Serializable{Typ: "integer", Num: -2}
	}
	stringMU.RLock()
	defer stringMU.RUnlock()

	val, ok := stringData[args[0].Bulk]
	if !ok || !val.will_expire {
		return resp.Serializable{Typ: "integer", Num: -1}
	}

	return resp.Serializable{Typ: "integer", Num: int(time.Until(val.expiration).Milliseconds())}
}

func expire(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "EXPIRE"}.Error()}
	}
	str, ok := getString(args[0].Bulk)
	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	expr_seconds, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "EXPIRE"}.Error()}
	}
	stringMU.Lock()
	stringData[args[0].Bulk] = stringDataItem{str: str, will_expire: true, expiration: time.Now().Add(time.Second * time.Duration(expr_seconds))}
	stringMU.Unlock()
	return resp.Serializable{Typ: "integer", Num: 1}

}
