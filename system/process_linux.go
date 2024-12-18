//go:build linux
// +build linux

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"v2k.io/gox/strings"
)

func getProcessById(pid int) (*Process, error) {
	// 执行文件
	exePath := fmt.Sprintf("/proc/%d/exe", pid)
	path, err := os.Readlink(exePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read exe link: %w", err)
	}

	dir := filepath.Dir(path)
	exe := filepath.Base(path)

	// 执行命令
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(data), "\x00", " ")
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("no command line information found")
	}
	command := parts[0]
	args := parts[1:]

	process := &Process{
		BinPath: dir,
		Exe:     exe,
		Command: command,
		Args:    args,
	}
	return process, nil
}
