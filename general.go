package main

func ping(args []Serializable) Serializable {
	return Serializable{typ: "bulk", bulk: "PONG"}
}

func lolwut(args []Serializable) Serializable {
	return Serializable{typ: "bulk", bulk: "GoKV 0.1 :):):)\r\n"}
}

func flushall(args []Serializable) Serializable {
	setMU.Lock()
	setData = map[string]string{}
	setMU.Unlock()
	HsetMU.Lock()
	HsetData = map[string]map[string]string{}
	HsetMU.Unlock()
	return Serializable{typ: "string", str: "OK"}
}
