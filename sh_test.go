package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	editor "mrqzzz/kpad/editor"
	"testing"
)

func TestExplain(t *testing.T) {
	//t.SkipNow()
	explain, err := editor.ExecKubectlExplain("pod.spec")
	assert.NoError(t, err)
	fmt.Println(string(explain))
	tree := editor.BuildExplainFieldsTree(explain)
	fmt.Println(tree)
}

func TestGetLeftmostWordAtLine(t *testing.T) {
	txt := `row0:
  row1:                
    row2: #hello
3
      row4:
  row5: |
    ### Heading

    * Bullet
    
  row10:
11
    12
🙂
x`

	e := editor.NewEditor(100, 100)
	e.StringToBuf(txt)
	st, x1, x2 := editor.GetLeftmostWordAtLine(e.Buf[0])
	assert.Equal(t, "row0:", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 5, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[1])
	assert.Equal(t, "row1:", st)
	assert.Equal(t, 2, x1)
	assert.Equal(t, 7, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[2])
	assert.Equal(t, "row2:", st)
	assert.Equal(t, 4, x1)
	assert.Equal(t, 9, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[3])
	assert.Equal(t, "3", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 1, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[4])
	assert.Equal(t, "row4:", st)
	assert.Equal(t, 6, x1)
	assert.Equal(t, 11, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[5])
	assert.Equal(t, "row5:", st)
	assert.Equal(t, 2, x1)
	assert.Equal(t, 7, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[7])
	assert.Equal(t, "", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 0, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[8])
	assert.Equal(t, "*", st)
	assert.Equal(t, 4, x1)
	assert.Equal(t, 5, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[9])
	assert.Equal(t, "", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 0, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[11])
	assert.Equal(t, "11", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 2, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[13])
	assert.Equal(t, "🙂", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 1, x2)

	st, x1, x2 = editor.GetLeftmostWordAtLine(e.Buf[14])
	assert.Equal(t, "x", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 1, x2)

}

func TestBuildCurrentPath(t *testing.T) {
	txt := `apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:

  - name: nginx


    image: nginx:1.14.2
    ports:
    - containerPort: 80
        
         
`
	e := editor.NewEditor(100, 100)
	e.StringToBuf(txt)
	st, _ := editor.BuildCurrentPath(e, 6, 12)
	assert.Equal(t, "spec.containers.ports", st)

}

func TestBuildCurrentPath2(t *testing.T) {
	txt := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.16.1 # Update the version of nginx from 1.14.2 to 1.16.1
        ports:
        - containerPort: 80
        
         
`
	e := editor.NewEditor(100, 100)
	e.StringToBuf(txt)

	st, _ := editor.BuildCurrentPath(e, 3, 2)
	assert.Equal(t, "metadata", st)

	st, _ = editor.BuildCurrentPath(e, 5, 7)
	assert.Equal(t, "spec.selector.matchLabels", st)

	st, _ = editor.BuildCurrentPath(e, 11, 18)
	assert.Equal(t, "spec.template.spec.containers.ports.containerPort", st)

	st, _ = editor.BuildCurrentPath(e, 2, 1)
	assert.Equal(t, "kind", st)

	st, _ = editor.BuildCurrentPath(e, 1, 7)
	assert.Equal(t, "spec", st)

}
