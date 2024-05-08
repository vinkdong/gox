package scram_sha_256

import (
	"fmt"
	"testing"
)

func TestParseSCRAMMessage(t *testing.T) {
	msg1 := " n,,n=,r=X/FS0mF/ARuLyuiCOaG2YgGG"
	p1, err := ParseSCRAMMessage(msg1)
	if err != nil {
		t.Fatal(err)
	}
	msg2 := "!n,,n=*,r=X/FS0mF/ARuLyuiCOaG2YgGG"
	p2, err := ParseSCRAMMessage(msg2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(p1, p2)
}
