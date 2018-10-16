package log

import "testing"

func TestInfo(t *testing.T) {
	Info("this is info")
}

func TestError(t *testing.T) {
	Error("this is error")
}

func TestSuccess(t *testing.T)  {
	Success("this is success")
}
