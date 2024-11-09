package utils

import "strconv"

func IpPortJoin(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}

func IsEmpty(s string) bool {
	return &s == nil || len(s) == 0 || s == "" || s == "\n"
}
