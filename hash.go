package main

func hset(args []Serializable) Serializable {

	if len(args)%2 != 1 || len(args) < 3 {
		return Serializable{typ: "error", str: "incorrect number of arguements - must have hash key and then key value pairs"}
	}
	HsetMU.Lock()
	HsetData[args[0].bulk] = map[string]string{}
	for i := 1; i < len(args); i += 2 {
		HsetData[args[0].bulk][args[i].bulk] = args[i+1].bulk
	}
	HsetMU.Unlock()

	return Serializable{typ: "integer", num: (len(args) - 1) / 2}
}

func hget(args []Serializable) Serializable {
	if len(args) < 2 {
		return Serializable{typ: "error", str: "incorrect number of arguements - must have hash key and then value key"}
	}

	HsetMU.RLock()
	defer HsetMU.RUnlock()

	_, ok := HsetData[args[0].bulk]

	if !ok {
		return Serializable{typ: "null"}
	}

	var resultList []Serializable
	for i := 1; i < len(args); i++ {
		value, ok := HsetData[args[0].bulk][args[i].bulk]
		if !ok {
			resultList = append(resultList, Serializable{typ: "null"})
			continue
		}
		resultList = append(resultList, Serializable{typ: "bulk", bulk: value})

	}
	return Serializable{typ: "array", array: resultList}

}

func hexists(args []Serializable) Serializable {
	if len(args) != 2 {
		return Serializable{typ: "error", str: "incorrect number of args, expected hashkey and field key"}
	}
	HsetMU.RLock()
	_, okKey := HsetData[args[0].bulk]
	_, okValue := HsetData[args[0].bulk][args[1].bulk]
	HsetMU.RUnlock()

	if !okKey || !okValue {
		return Serializable{typ: "integer", num: 0}
	}

	return Serializable{typ: "integer", num: 1}
}

func hstrlen(args []Serializable) Serializable {
	if len(args) != 2 {
		return Serializable{typ: "error", str: "incorrect number of args"}
	}
	HsetMU.RLock()
	val, ok := HsetData[args[0].bulk][args[1].bulk]
	HsetMU.RUnlock()
	if !ok {
		return Serializable{typ: "integer", num: 0}
	}
	return Serializable{typ: "integer", num: len(val)}
}
