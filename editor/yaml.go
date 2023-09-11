package editor

import (
	"errors"
	"strings"
)

type YNode struct {
	Key      string
	Value    string
	Children []*YNode
}

func ReadFlatYaml(s string) (root *YNode, err error) {
	lines := strings.Split(s, "\n")
	root = &YNode{Children: []*YNode{}}
	for _, line := range lines {
		lineTrimmed := strings.TrimSpace(line)
		if lineTrimmed == "" || strings.HasPrefix(lineTrimmed, "#") {
			continue
		}
		sts := strings.Split(line, ":")
		if len(sts) < 2 {
			return nil, errors.New("Ivalid config. It must be a flat 'key: value' list")
		}
		sts[1] = strings.Join(sts[1:], ":")
		root.Children = append(root.Children, &YNode{
			Key:      strings.TrimSpace(sts[0]),
			Value:    strings.TrimSpace(sts[1]),
			Children: []*YNode{},
		})

	}
	return root, nil
}
