package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExplain(t *testing.T) {
	explain, err := ExecExplain("pod.spec")
	assert.NoError(t, err)
	fmt.Println(string(explain))
	tree := BuildExplainFieldsTree(explain)
	fmt.Println(tree)
}

func TestGetWordAtPos(t *testing.T) {
	txt := `row0:
  row1:
    row2:
    row3:
      row4:
`
	st, x1, x2 := getWordAtPos(txt, 1, 0)
	assert.Equal(t, st, "row0:")
	assert.Equal(t, x1, 0)
	assert.Equal(t, x2, 5)

	st, x1, x2 = getWordAtPos(txt, 3, 1)
	assert.Equal(t, st, "row1:")
	assert.Equal(t, x1, 2)
	assert.Equal(t, x2, 7)

}
