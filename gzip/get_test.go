package gzip

import (
	"fmt"
	"testing"
)

func TestGetFileFromGzip(t *testing.T) {
	data, err := GetFileFromGzip("/Users/vink/tmp/ylip-1.34-charts/zknow-platform-0.26.0.tgz", "values-template.yml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}
