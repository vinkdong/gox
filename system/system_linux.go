//go:build linux
// +build linux

package system

import (
	"os/exec"
	"strings"
)

// 获取 Linux 系统 UUID
func getSystemUUID() (string, error) {
	cmd := exec.Command("cat", "/sys/class/dmi/id/product_uuid")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
