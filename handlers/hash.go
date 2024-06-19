package handlers

import "github.com/idugan100/GoKV/resp"

func hset(args []resp.Serializable) resp.Serializable {

	if len(args)%2 != 1 || len(args) < 3 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HSET"}.Error()}
	}
	hsetMU.Lock()
	hsetData[args[0].Bulk] = map[string]string{}
	for i := 1; i < len(args); i += 2 {
		hsetData[args[0].Bulk][args[i].Bulk] = args[i+1].Bulk
	}
	hsetMU.Unlock()

	return resp.Serializable{Typ: "integer", Num: (len(args) - 1) / 2}
}

func hget(args []resp.Serializable) resp.Serializable {
	if len(args) < 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HGET"}.Error()}
	}

	hsetMU.RLock()
	defer hsetMU.RUnlock()

	_, ok := hsetData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	var resultList []resp.Serializable
	for i := 1; i < len(args); i++ {
		value, ok := hsetData[args[0].Bulk][args[i].Bulk]
		if !ok {
			resultList = append(resultList, resp.Serializable{Typ: "null"})
			continue
		}
		resultList = append(resultList, resp.Serializable{Typ: "bulk", Bulk: value})

	}
	return resp.Serializable{Typ: "array", Array: resultList}

}

func hexists(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HEXISTS"}.Error()}
	}
	hsetMU.RLock()
	_, okKey := hsetData[args[0].Bulk]
	_, okValue := hsetData[args[0].Bulk][args[1].Bulk]
	hsetMU.RUnlock()

	if !okKey || !okValue {
		return resp.Serializable{Typ: "integer", Num: 0}
	}

	return resp.Serializable{Typ: "integer", Num: 1}
}

func hstrlen(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HSTRLEN"}.Error()}
	}
	hsetMU.RLock()
	val, ok := hsetData[args[0].Bulk][args[1].Bulk]
	hsetMU.RUnlock()
	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	return resp.Serializable{Typ: "integer", Num: len(val)}
}

func hlen(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HLEN"}.Error()}
	}

	hsetMU.RLock()
	val, ok := hsetData[args[0].Bulk]
	hsetMU.RUnlock()

	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	return resp.Serializable{Typ: "integer", Num: len(val)}

}

func hgetall(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HGETALL"}.Error()}
	}
	hsetMU.RLock()
	val, ok := hsetData[args[0].Bulk]
	hsetMU.RUnlock()
	if !ok {
		return resp.Serializable{Typ: "array"}
	}
	var results []resp.Serializable
	for key := range val {
		results = append(results, resp.Serializable{Typ: "bulk", Bulk: val[key]})
	}
	return resp.Serializable{Typ: "array", Array: results}

}

func hsetnx(args []resp.Serializable) resp.Serializable {
	if len(args) != 3 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HSETNX"}.Error()}
	}
	hsetMU.Lock()
	defer hsetMU.Unlock()

	_, ok := hsetData[args[0].Bulk][args[1].Bulk]

	if ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}

	hsetData[args[0].Bulk][args[1].Bulk] = args[2].Bulk

	return resp.Serializable{Typ: "integer", Num: 1}
}

func hdel(args []resp.Serializable) resp.Serializable {
	if len(args) < 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "HDEL"}.Error()}
	}
	numFields := len(args) - 1
	numDeleted := 0

	hsetMU.Lock()
	defer hsetMU.Unlock()
	for i := range numFields {
		if _, ok := hsetData[args[0].Bulk][args[i+1].Bulk]; ok {
			delete(hsetData[args[0].Bulk], args[i+1].Bulk)
			numDeleted++
		}
	}

	return resp.Serializable{Typ: "integer", Num: numDeleted}

}
