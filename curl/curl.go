package curl

import (
	"github.com/imroc/req"
	shellwords "github.com/mattn/go-shellwords"
	"strings"
)

type Curl struct {
	Url     string            `json:"host"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Data    string            `json:"data"`
}

func (c *Curl) Call() (*req.Resp, error) {

	header := req.Header{}
	for h, v := range c.Headers {
		header[h] = v
	}
	return req.Do(c.Method, c.Url, header, c.Data)
}

func ParseCurlCommand(curlCmd string) (*Curl, error) {
	args, err := shellwords.Parse(curlCmd)
	if err != nil {
		return nil, err
	}
	//re := regexp.MustCompile(`(?:'|")(.*?)(?:'|")|\S+`)
	//args := re.FindAllString(curlCmd, -1)

	method := ""
	url := ""
	headers := make(map[string]string)
	var data string

	// i=0是curl
	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch strings.Trim(arg, "\n") {
		case "-X", "--request":
			i++
			method = args[i]
		case "-H", "--header":
			i++
			headerParts := strings.SplitN(args[i], ": ", 2)
			headers[headerParts[0]] = headerParts[1]
		case "-d", "--data-raw", "--data":
			i++
			data = args[i]
		default:
			if !strings.HasPrefix(arg, "-") && url == "" {
				url = arg
			}
		}
	}

	if method == "" {
		if data != "" {
			method = "POST"
		} else {
			method = "GET"
		}
	}

	c := Curl{
		Url:     url,
		Method:  method,
		Headers: headers,
		Data:    data,
	}
	return &c, nil
}
