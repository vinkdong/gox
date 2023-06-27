package curl

import (
	shellwords "github.com/mattn/go-shellwords"
	"strings"
)

type Curl struct {
	Url     string            `json:"host"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Data    string            `json:"data"`
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
		switch arg {
		case "-X":
			i++
			method = args[i]
		case "-H":
			i++
			headerParts := strings.SplitN(args[i], ": ", 2)
			headers[headerParts[0]] = headerParts[1]
		case "-d", "--data-raw":
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
