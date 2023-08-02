package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExplain(t *testing.T) {
	t.SkipNow()
	explain, err := ExecExplain("pod.spec")
	assert.NoError(t, err)
	fmt.Println(string(explain))
	tree := BuildExplainFieldsTree(explain)
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

	st, x1, x2 := getLeftmostWordAtLine(txt, 0)
	assert.Equal(t, "row0:", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 5, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 1)
	assert.Equal(t, "row1:", st)
	assert.Equal(t, 2, x1)
	assert.Equal(t, 7, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 2)
	assert.Equal(t, "row2:", st)
	assert.Equal(t, 4, x1)
	assert.Equal(t, 9, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 3)
	assert.Equal(t, "3", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 1, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 4)
	assert.Equal(t, "row4:", st)
	assert.Equal(t, 6, x1)
	assert.Equal(t, 11, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 5)
	assert.Equal(t, "row5:", st)
	assert.Equal(t, 2, x1)
	assert.Equal(t, 7, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 7)
	assert.Equal(t, "", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 0, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 8)
	assert.Equal(t, "*", st)
	assert.Equal(t, 4, x1)
	assert.Equal(t, 5, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 9)
	assert.Equal(t, "", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 0, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 11)
	assert.Equal(t, "11", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 2, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 13)
	assert.Equal(t, "🙂", st)
	assert.Equal(t, 0, x1)
	assert.Equal(t, 4, x2)

	st, x1, x2 = getLeftmostWordAtLine(txt, 14)
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

	st := buildCurrentPath(txt, 6, 12)
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

	st := buildCurrentPath(txt, 3, 3)
	assert.Equal(t, "metadata", st)

	st = buildCurrentPath(txt, 5, 7)
	assert.Equal(t, "spec.selector.matchLabels", st)

	st = buildCurrentPath(txt, 9, 18)
	assert.Equal(t, "spec.template.spec.containers.ports", st)

	st = buildCurrentPath(txt, 2, 1)
	assert.Equal(t, "", st)

	st = buildCurrentPath(txt, 1, 7)
	assert.Equal(t, "spec", st)

}
