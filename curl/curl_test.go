package curl

import (
	"fmt"
	"testing"
)

func TestParseCurlCommand(t *testing.T) {
	curlTxt := `
curl 'https://api.example.com/intelligent/v1/228549383619211264/conversion/fixMessageText' -H 'accept: application/json' \
  -H 'authorization: bearer 4af03ecf-xxx-4c83-0000-3fa29657af7e' \
  -H 'content-type: application/json' \
  -H "x-domain-path: undefined" \
  -H 'x-language: undefined' \
  -H 'x-tenant-id: 228549383619210000' \
  -d '[200049383619211264,300001059980095488,400052212098428928,207965165382130000]'
`
	curl, err := ParseCurlCommand(curlTxt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(curl.Url)
}
