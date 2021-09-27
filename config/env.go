package config

import (
	"os"
	"strconv"
	"strings"
)

// Env 获取环境配置字符串
func Env(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return value
}

// EnvInt 获取环境配置数值
func EnvInt(key string, defaultValue ...int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return value
}

// EnvBool 获取环境配置布尔值
func EnvBool(key string, defaultValue ...bool) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return value
}
