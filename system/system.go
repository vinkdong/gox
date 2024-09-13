package system

import "v2k.io/gox/log"

func GetSystemUUID() string {
	uuid, err := getSystemUUID()
	if err != nil {
		log.ErrorLine(err)
	}
	return uuid
}
