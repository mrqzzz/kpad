package main

import (
	"fmt"
	"testing"
)

func TestExplain(t *testing.T) {
	explain, err := ExecExplain("pod.spec")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(explain))
	tree := BuildExplainFieldsTree(explain)
	fmt.Println(tree)
}
