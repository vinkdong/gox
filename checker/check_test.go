package checker

import (
	"testing"
	"os"
	"text/template"
)

type X struct {
	Status string

}

func TestChecker_CheckJson(t *testing.T) {
	c, _ := New("t0").Parse(`{{ if eq $.Status "red" }} T1 {{end}}`)
	x := &X{Status: "green"}
	c.execute(os.Stdout, x)

	tx, _ := template.New("t1").Parse(`{{ if eq $.Status "red" }} T1 {{end}}`)
	tx.Execute(os.Stdout,x)
	//fmt.Println(c)
}