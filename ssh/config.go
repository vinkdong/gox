package ssh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/vinkdong/gox/log"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config []HostConfig

type HostConfig struct {
	Host                  string
	ServerAliveInterval   int
	HostName              string
	User                  string
	Port                  string
	StrictHostKeyChecking string
	UserKnownHostsFile    string
	IdentityFile          string
	LogLevel              string
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
		fmt.Println(host)
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
