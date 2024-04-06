package utils

import "strconv"

func UrlCombine(ip string, port int, resource string) string {
	if resource != "" {
		return ip + ":" + strconv.Itoa(port) + "/" + resource
	} else {
		return ip + ":" + strconv.Itoa(port)
	}
}
