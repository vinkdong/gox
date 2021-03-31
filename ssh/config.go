package ssh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"v2k.io/gox/log"
)

// ssh config文件
type Config []HostConfig

type HostConfig struct {
	Host                  string
	ServerAliveInterval   int
	HostName              string
	User                  string
	Port                  int
	StrictHostKeyChecking string
	UserKnownHostsFile    string
	IdentityFile          string
	LogLevel              string
}

func (hc *HostConfig) Map(lambda func(key string, value interface{})) {
	s := reflect.ValueOf(hc).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		lambda(typeOfT.Field(i).Name, f.Interface())
	}
}

func (c *Config) SetHost(host HostConfig) {
	newConfig := make(Config, 0, len(*c)+1)
	isExist := false
	for _, v := range *c {
		if v.Host == host.Host {
			// for no change order
			newConfig = append(newConfig, host)
			isExist = true
			continue
		}
		newConfig = append(newConfig, v)
	}
	if !isExist {
		newConfig = append(newConfig, host)
	}
	*c = newConfig
}

func (c *Config) RemoveHost(host string) {
	newConfig := make(Config, 0, len(*c)+1)
	for _, v := range *c {
		if v.Host == host {
			continue
		}
		newConfig = append(newConfig, v)
	}
	*c = newConfig
}

func (c *Config) GetHost(host string) HostConfig {
	for _, v := range *c {
		if v.Host == host {
			return v
		}
	}
	return HostConfig{}
}

func (c *Config) Marshal() []byte {
	buf := bytes.Buffer{}

	l := func(key string, value interface{}) {
		if key == "Host" {
			return
		}
		switch value.(type) {
		case int:
			// 不输出为0的数据
			if value.(int) == 0 {
				return
			}
		case string:
			if value.(string) == "" {
				return
			}
		}
		buf.WriteString(fmt.Sprintf("   %s %v\n", key, value))
	}
	for _, v := range *c {
		buf.WriteString(fmt.Sprintf("Host %s\n", v.Host))
		v.Map(l)
	}
	return buf.Bytes()
}

func Unmarshal(cfg *Config, data []byte) {
	r := &bytes.Buffer{}
	r.Write(data)
	*cfg = *Scan(r)
}

func Read(file *os.File) *Config {
	return Scan(file)
}

func Scan(r io.Reader) *Config {
	s := bufio.NewScanner(r)
	c := make(Config, 0)
	var host HostConfig
	n := 0
	for s.Scan() {
		n++
		key, value := convertLine(s.Bytes())
		if key == "" {
			continue
		}
		if key == "Host" {
			if host.Host != "" {
				c = append(c, host)
			}
			host = HostConfig{}
		}
		if err := setHostConfig(&host, key, value); err != nil {
			log.Error(err)
		}
	}
	if &host != nil {
		c = append(c, host)
	}
	return &c
}

// convert line to key:value
//    Host   abc
func convertLine(lineBytes []byte) (string, string) {
	key := bytes.Buffer{}
	value := bytes.Buffer{}
	readKey := false
	readValue := false
	for _, v := range lineBytes {
		if v == ' ' {
			if key.Len() == 0 {
				continue
			} else {
				readKey = false
			}
		}
		if v != ' ' {
			if key.Len() == 0 {
				readKey = true
			}
			if readKey {
				key.WriteByte(v)
				continue
			} else {
				readValue = true
			}
		}
		if readValue {
			value.WriteByte(v)
		}
	}
	return key.String(), strings.Trim(value.String(), " ")
}

// todo: should no use reflect
func setHostConfig(config *HostConfig, field string, value string) error {
	v := reflect.ValueOf(config).Elem().FieldByName(field)
	if !v.IsValid() {
		return errors.New(fmt.Sprintf("no such filed : %s in HostConfig Type", field))
	}
	switch v.Type().Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		v.SetInt(int64(i))
	}
	return nil
}
