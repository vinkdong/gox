//go:build darwin
// +build darwin

package system

import (
	"fmt"
	"os/exec"
	"regexp"
)

// 获取 macOS 系统 UUID
func getSystemUUID() (string, error) {
	cmd := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 正则表达式匹配 UUID
	re := regexp.MustCompile(`"IOPlatformUUID"\s*=\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("unable to find system UUID on macOS")
}
