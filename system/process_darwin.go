//go:build darwin
// +build darwin

package system

import "errors"

func getProcessById(pid int) (*Process, error) {
	return nil, errors.New("not Implement")
}
