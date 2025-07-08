package slice

import (
	"fmt"
	"testing"
	"time"
)

func TestDifference(t *testing.T) {
	sla := []string{"a", "b", "c", "f"}
	slb := []string{"b", "c", "d", "e"}
	diff := Difference(sla, slb)
	fmt.Println(diff)
	if len(diff) != 2 {
		t.Error("slice difference set error")
	}
	diff2 := Difference(diff, sla)
	if len(diff2) != 0 {
		t.Error("slice difference set error")
	}
}

func TestUnionString(t *testing.T) {
	s0 := []string{"a", "b", "d"}
	s1 := []string{"a", "d", "e"}
	s2 := []string{"d", "x", "d"}
	var s3 []string
	s4 := UnionString(s0, s1, s2, s3)
	if len(s4) != 5 {
		t.Fatal("slice union len should be 5 ")
	}
}

func TestStructToSortedFlatSlice(t *testing.T) {
	type C struct {
		A string         `json:"ca"`
		B *time.Duration `json:"cb"`
	}
	type Cfg struct {
		A string `json:"aa"`
		B int    `json:"ab"`
		C *C     `json:"ac"`
	}
	_t := time.Second * 14
	c := Cfg{
		A: "abc",
		B: 5,
		C: &C{
			A: "a",
			B: &_t,
		},
	}
	var flatSlice []KeyValue
	StructToSortedFlatSlice("", c, &flatSlice)
}
