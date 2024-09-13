//go:build windows
// +build windows

package system

import (
	"golang.org/x/sys/windows/registry"
)

// 获取 Windows 系统 UUID
func getSystemUUID() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Cryptography`, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	uuid, _, err := k.GetStringValue("MachineGuid")
	if err != nil {
		return "", err
	}
	return uuid, nil
}
