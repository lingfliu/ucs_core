package utils

import "strconv"

func IpPortJoin(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func Array2String(array []any) string {
	str := ""
	for _, v := range array {
		str += string(v)
		str += " "
	}
	return str[:len(str)-1]
}
