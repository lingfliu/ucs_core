package utils

import "strconv"

func IpPortJoin(ip string, port int) string {
	return ip + ":" + strconv.Itoa(port)
}
