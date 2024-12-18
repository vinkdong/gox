package system

type Process struct {
	BinPath string
	Command string
	Exe     string
	Args    []string
}

func GetProcessInfo(pid int) (*Process, error) {
	return getProcessById(pid)
}
