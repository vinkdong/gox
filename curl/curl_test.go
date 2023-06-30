package curl

import (
	"fmt"
	"testing"
)

func TestParseCurlCommand(t *testing.T) {
	curlTxt := `
curl --location --request PUT 'https://api.example.com/intelligent/v1/0/bots/flow/fix_1_22_search' \
--header 'authorization: bearer fa956d2f-4419-4f83-b37a-6303ca565ce0' \
--header 'Content-Type: application/json' \
--data '{
    "exceptBotIds": [],
    "exceptTenantIds": []
}'
`
	curl, err := ParseCurlCommand(curlTxt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(curl.Url)
}
