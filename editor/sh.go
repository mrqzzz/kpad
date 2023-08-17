package editor

import (
	"os/exec"
	"strings"
)

type Node struct {
	FieldName string
	FieldType string
	Indent    int
	Parent    *Node
	Children  []*Node
}

func ExecKubectlExplain(path string) ([]byte, error) {
	args := []string{"explain", "--recursive", path}
	out, err := exec.Command("kubectl", args...).Output()
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

func ExecKubectlApiResources() ([]byte, error) {
	args := []string{"api-resources"}
	out, err := exec.Command("kubectl", args...).Output()
	return out, err
}

func BuildApiResourcesList(bytes []byte) (names []string, versions []string) {
	names = []string{}
	versions = []string{}
	sts := strings.Split(string(bytes), "\n")
	if len(sts) > 0 {
		idxNameFrom := strings.Index("NAME", sts[0])
		idxNameTo := strings.Index("SHORTNAMES", sts[0])
		idxAPiVersionFrom := strings.Index("APIVERSION", sts[0])
		idxAPiVersionTo := strings.Index("NAMESPACED", sts[0])
		for i := 1; i < len(sts); i++ {
			names = append(names, strings.Trim(sts[i][idxNameFrom:idxNameTo], " "))
			versions = append(versions, strings.Trim(sts[i][idxAPiVersionFrom:idxAPiVersionTo], " "))
		}
	}
	return
}
