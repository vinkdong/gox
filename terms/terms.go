package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"v2k.io/gox/log"
)

type Input struct {
	//是否启用
	Enabled bool
	//需要正则匹配
	Regex string
	//提示信息
	Tip string
	//是否是密码输入
	Password bool
	//输入必须包含
	Include []KV
	//输入不能包含
	Exclude []KV
	//是否需要两次输入
	Twice   bool
	Default string
}

type KV struct {
	Name  string
	Value string
}

func (input *Input) AskBool() bool {
	in, err := input.AcceptUserInput()
	if err != nil {
		fmt.Println("输入错误")
		return input.AskBool()
	}
	if in == "" || strings.Trim(in, " ") == "" {
		in = input.Default
	}
	switch in {
	case "yes", "y", "Y", "Yes", "true", "t", "TRUE", "1":
		return true
	case "no", "n", "N", "false", "f", "FALSE", "0":
		return false
	default:
		fmt.Println("输入不满足需求")
		return input.AskBool()
	}
}

func (input *Input) AskInt() int {
	in, err := input.AcceptUserInput()
	if err != nil {
		fmt.Println("输入错误")
		return input.AskInt()
	}
	if in == "" || strings.Trim(in, " ") == "" {
		in = input.Default
	}
	i, err := strconv.Atoi(in)
	if err != nil {
		fmt.Println("请输入整数")
		return input.AskInt()
	}
	return i
}

func (input *Input) AskString() string {
	in, err := input.AcceptUserInput()
	if err != nil {
		fmt.Println("输入错误")
		return input.AskString()
	}
	if in == "" || strings.Trim(in, " ") == "" {
		in = input.Default
	}
	return in
}

func (input *Input) AcceptUserInput() (string, error) {
	if input.Password {
		return input.AcceptUserPassword()
	}
start:
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(input.Tip + " [" + input.Default + "]:  ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.Trim(text, "\n")
	if text == "" {
		text = input.Default
	}

	if !input.CheckMatch(text) {
		goto start
	}
	return text, nil
}

func (input *Input) AcceptUserPassword() (string, error) {
start:
	fmt.Print(input.Tip)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if !input.CheckMatch(string(bytePassword[:])) {
		goto start
	}

	if !input.Twice {
		return string(bytePassword[:]), nil
	}

	fmt.Print("请再输入一次:")
	bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	if len(bytePassword2) != len(bytePassword) {
		log.Error("两次输入长度不符")
		goto start
	}
	for k, v := range bytePassword {
		if bytePassword2[k] != v {
			log.Error("两次输入不同")
			goto start
		}
	}

	fmt.Println("waiting...")
	return string(bytePassword[:]), nil
}

func (input *Input) CheckMatch(value string) bool {

	r := regexp.MustCompile(input.Regex)
	if !r.MatchString(value) {
		log.Errorf("输入不满足需求")
		return false
	}

	for _, include := range input.Include {
		r := regexp.MustCompile(include.Value)
		if !r.MatchString(value) {
			log.Errorf(include.Name)
			return false
		}
	}

	for _, exclude := range input.Exclude {
		r := regexp.MustCompile(exclude.Value)
		if r.MatchString(value) {
			log.Errorf(exclude.Name)
			return false
		}
	}
	return true
}
