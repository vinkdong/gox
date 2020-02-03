package downloader

import (
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	statusCode, err := Download("https://vinkdong.com/sis", "a/b/c/sxs")
	if err != nil {
		t.Fatal(err)
	}
	if statusCode != 404 {
		t.Fatal("status code should be 404")
	}
	os.RemoveAll("a")
}
