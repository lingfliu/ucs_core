package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GenMqttCliId() string {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")
	return "mm-" + id
}

func CheckMqttServerUrl(url string) bool {
	if url == "" {
		return false
	}

	if !(strings.HasPrefix(url, "mqtt://") || strings.HasPrefix(url, "mqtts://") || strings.HasPrefix(url, "ws://") || strings.HasPrefix(url, "wss://")) {
		return false
	}

	return true
}
