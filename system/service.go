package system

import (
	"bytes"
	"os/exec"
	"strings"
)

type Service struct {
	Name        string
	ExecStart   string
	ExecStop    string
	Environment []string
}

func ListServiceByPrefix(prefix string) ([]string, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-pager")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	lines := strings.Split(out.String(), "\n")
	var redisServices []string
	for _, line := range lines {
		if strings.HasPrefix(line, "redis") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				serviceName := fields[0]
				redisServices = append(redisServices, serviceName)
			}
		}
	}
	return redisServices, nil
}

func GetService(name string) (*Service, error) {
	data, err := GetUnitFileContent(name)
	if err != nil {
		return nil, err
	}
	s := parseUnitFile(data)
	s.Name = name
	return s, nil
}

func GetUnitFileContent(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "cat", serviceName)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func parseUnitFile(content string) *Service {
	var service Service
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if keyValue := strings.SplitN(line, "=", 2); len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]

			switch key {
			case "ExecStart":
				service.ExecStart = value
			case "ExecStop":
				service.ExecStop = value
			case "Environment":
				service.Environment = append(service.Environment, value)
			}
		}
	}
	return &service
}
