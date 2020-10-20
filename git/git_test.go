package git

import (
	"fmt"
	"testing"
)

func TestExtractRepositoryPath(t *testing.T) {
	t1 := "https://github.com/vinkdong/records.git"
	r1 := ExtractRepositoryPath(t1)

	if r1 != "vinkdong/records" {
		t.Fatal("提取失败")
	}

	t2 := "https://someuser:password:github.com/vinkdong/records.git"
	r2 := ExtractRepositoryPath(t2)
	if r2 != "vinkdong/records" {
		t.Fatal("提取失败")
	}

	t3 := "git@github.com:vinkdong/records.git"
	r3 := ExtractRepositoryPath(t3)
	if r3 != "vinkdong/records" {
		t.Fatal("提取失败")
		fmt.Println(r3)
	}
}
