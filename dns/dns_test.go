package dns

import (
	"fmt"
	"testing"
)

func TestDNS_LookUpNS(t *testing.T) {
	d, err := New("114.114.114.114")
	fmt.Println(err)
	fmt.Println(d.LookupIPAddr("vinkdong.com"))
}
