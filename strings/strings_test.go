package strings

import (
	"fmt"
	"testing"
)

func TestSplit(t *testing.T) {
	s := Split("abc\x00de fh", "\x00", " ")
	fmt.Println(s)
}
