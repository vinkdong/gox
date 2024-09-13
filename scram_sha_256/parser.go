package scram_sha_256

import (
	"fmt"
	"strconv"
	"strings"
)

type SCRAMParams struct {
	DefaultUsername string // 没有等号的默认用户名
	Username        string // 实际用户名列表
	ClientNonce     string // 客户端随机数
	Salt            string // 盐值
	Iterations      int    // 迭代次数
}

// ParseSCRAMMessage 解析 SCRAM 消息并返回一个 SCRAMParams 结构
func ParseSCRAMMessage(message string) (*SCRAMParams, error) {
	params := &SCRAMParams{}
	message = strings.TrimPrefix(message, " ")
	message = strings.TrimPrefix(message, "!")
	pairs := strings.Split(message, ",") // 分割字符串得到各个键值对
	for _, pair := range pairs {
		if pair == "" {
			continue // 忽略空字符串，防止干扰解析
		}
		kv := strings.SplitN(pair, "=", 2) // 用等号分割键和值
		key := kv[0]
		var value string
		if len(kv) >= 2 {
			value = kv[1]
		}
		switch key {
		case "n":
			if len(kv) == 1 {
				params.DefaultUsername = kv[0]
				continue
			}
			params.Username = value
		case "r":
			params.ClientNonce = value
		case "s":
			params.Salt = value
		case "i":
			iteration, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			params.Iterations = iteration
		default:
			fmt.Printf("Unrecognized key: %s\n", key)
		}
	}
	return params, nil
}
