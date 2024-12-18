package system

import (
	"fmt"
	"testing"
)

func TestListServiceByPrefix(t *testing.T) {
	ss, err := ListServiceByPrefix("redi")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ss)
	s, err := GetService(ss[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
