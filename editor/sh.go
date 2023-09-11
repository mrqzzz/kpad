package editor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type Node struct {
	FieldName string
	FieldType string
	Indent    int
	Parent    *Node
	Children  []*Node
}

func (e *Editor) ExecKubectlExplain(path string) ([]byte, error) {
	args := []string{"explain", "--recursive", path}
	e.StatusBar.DrawInfo(fmt.Sprint(KUBECTL, args))
	out, err := exec.Command(KUBECTL, args...).Output()
	return out, err
}

func BuildExplainFieldsTree(explain []byte) *Node {
	sts := strings.Split(string(explain), "\n")
	var root *Node = nil
	startRow := -1
	for i, st := range sts {
		if st == "FIELDS:" {
			startRow = i
			break
		}
	}
	if startRow == -1 {
		return nil
	}
	var parent *Node = nil
	var prevNode *Node = nil
	currentInd := -1
	for i := startRow; i < len(sts); i++ {
		if sts[i] == "" {
			continue
		}
		ind, tabPos := getIndentAndTabPos(sts[i])
		n := &Node{
			FieldName: sts[i][ind:tabPos],
			FieldType: sts[i][tabPos+1:],
			Indent:    ind,
			Parent:    parent,
			Children:  []*Node{},
		}
		if ind > currentInd {
			if root == nil {
				root = n
			}
			n.Parent = prevNode
			parent = prevNode
		}
		if ind > currentInd || ind == currentInd {
			if parent != nil {
				parent.Children = append(parent.Children, n)
			}
		} else {
			sib := findPreviousSibling(n)
			if sib != nil && sib.Parent != nil {
				parent = sib.Parent
				n.Parent = parent
				parent.Children = append(sib.Parent.Children, n)
			}
		}
		prevNode = n
		currentInd = ind
	}
	return root
}

func findPreviousSibling(n *Node) *Node {
	parent := n.Parent
	for {
		if parent == nil {
			return nil
		}
		if parent.Indent == n.Indent {
			return parent
		}
		parent = parent.Parent
	}
}

func getIndentAndTabPos(st string) (indent int, tabPos int) {
	for i := 0; i < len(st); i++ {
		if st[i] == ' ' {
			indent++
		} else {
			break
		}
	}
	for i := indent; i < len(st); i++ {
		if st[i] == '\t' {
			tabPos = i
			break
		}
	}
	return
}

///////

func (e *Editor) ExecKubectlApiResources() ([]byte, error) {
	args := []string{"api-resources", "--sort-by=name"}
	e.StatusBar.DrawInfo(fmt.Sprint(KUBECTL, args))
	out, err := exec.Command(KUBECTL, args...).Output()
	return out, err
}

func BuildApiResourcesList(bytes []byte) (names []string, versions []string) {
	names = []string{}
	versions = []string{}
	sts := strings.Split(string(bytes), "\n")
	if len(sts) > 0 {
		idxNameFrom := strings.Index(sts[0], "NAME")
		idxNameTo := strings.Index(sts[0], "SHORTNAMES")
		idxAPiVersionFrom := strings.Index(sts[0], "APIVERSION")
		idxAPiVersionTo := strings.Index(sts[0], "NAMESPACED")
		for i := 1; i < len(sts); i++ {
			if len(sts[i]) > idxAPiVersionTo {
				names = append(names, strings.Trim(sts[i][idxNameFrom:idxNameTo], " "))
				versions = append(versions, strings.Trim(sts[i][idxAPiVersionFrom:idxAPiVersionTo], " "))
			}
		}
	}
	return
}

func IsWindows() bool {
	return strings.HasPrefix(runtime.GOOS, "windows")
}
