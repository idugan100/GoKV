package main

import "github.com/idugan100/GoKV/resp"

func hset(args []resp.Serializable) resp.Serializable {

	if len(args)%2 != 1 || len(args) < 3 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements - must have hash key and then key value pairs"}
	}
	HsetMU.Lock()
	HsetData[args[0].Bulk] = map[string]string{}
	for i := 1; i < len(args); i += 2 {
		HsetData[args[0].Bulk][args[i].Bulk] = args[i+1].Bulk
	}
	HsetMU.Unlock()

	return resp.Serializable{Typ: "integer", Num: (len(args) - 1) / 2}
}

func hget(args []resp.Serializable) resp.Serializable {
	if len(args) < 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of arguements - must have hash key and then value key"}
	}

	HsetMU.RLock()
	defer HsetMU.RUnlock()

	_, ok := HsetData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	var resultList []resp.Serializable
	for i := 1; i < len(args); i++ {
		value, ok := HsetData[args[0].Bulk][args[i].Bulk]
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
		return resp.Serializable{Typ: "error", Str: "incorrect number of args, expected hashkey and field key"}
	}
	HsetMU.RLock()
	_, okKey := HsetData[args[0].Bulk]
	_, okValue := HsetData[args[0].Bulk][args[1].Bulk]
	HsetMU.RUnlock()

	if !okKey || !okValue {
		return resp.Serializable{Typ: "integer", Num: 0}
	}

	return resp.Serializable{Typ: "integer", Num: 1}
}

func hstrlen(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: "incorrect number of args"}
	}
	HsetMU.RLock()
	val, ok := HsetData[args[0].Bulk][args[1].Bulk]
	HsetMU.RUnlock()
	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	return resp.Serializable{Typ: "integer", Num: len(val)}
}