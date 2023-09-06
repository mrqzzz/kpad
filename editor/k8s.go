package editor

import (
	"fmt"
	"strings"
)

//func populateDropdown(root *Node) {
//	items := []string{}
//	if root != nil {
//		for _, child := range root.Children {
//			items = append(items, child.FieldName+"   "+child.FieldType)
//		}
//
//	}
//	//drop.SetOptions(items, nil)
//}

func GenerateResourceTemplate(resourceName string, resourceVersion string) string {
	st := `apiVersion: %s
kind: %s
metadata:`
	return fmt.Sprintf(st, resourceVersion, resourceName)
}

// BuildCurrentPath returns the k8s yaml path ("deployment.spec.template.spec") from the cursor up
// also returns the possible same-level siblings.
// The passed x and y are relative to buf, not the cursor
func BuildCurrentPath(e *Editor, x int, y int) (stringPath string, existingSiblings map[string]interface{}) {
	var path []string
	var kind = ""
	xx := x
	existingSiblings = make(map[string]interface{})

	// build the path
	for i := y; i >= 0; i-- {
		fromIdx := e.findPrevLineFeed(i) + 1
		toIdx := e.findNextLineFeed(i)
		txt := runesJoin(e.Buf[fromIdx : toIdx+1])
		st, x1, _ := GetLeftmostWordAtLine(txt)
		if strings.HasPrefix(st, "---") {
			break
		}
		if kind == "" && st == "kind:" && x1 == 0 {
			spl := strings.Split(string(txt), ":")
			if len(spl) > 0 {
				kind = strings.Trim(spl[1], " \t\n") + "."
			}
		} else {
			if x1 < xx && st != "" {
				if st[len(st)-1:] == ":" {
					st = st[:len(st)-1]
				}
				path = append([]string{st}, path...)
				xx = x1
			}
		}
	}
	stringPath = kind + strings.Join(path, ".")

	// build the existingSiblings list, moving up
	for i := y; i >= 0; i-- {
		fromIdx := e.findPrevLineFeed(i) + 1
		toIdx := e.findNextLineFeed(i)
		txt := runesJoin(e.Buf[fromIdx : toIdx+1])
		st, x1, _ := GetLeftmostWordAtLine(txt)
		if x1 < x {
			break
		}
		if x1 == x && st != "" {
			if st[len(st)-1:] == ":" {
				st = st[:len(st)-1]
			}
			existingSiblings[st] = true
		}
	}

	// build the existingSiblings list, moving down
	for i := y + 1; i < len(e.Buf); i++ {
		fromIdx := e.findPrevLineFeed(i) + 1
		toIdx := e.findNextLineFeed(i)
		txt := runesJoin(e.Buf[fromIdx : toIdx+1])
		st, x1, _ := GetLeftmostWordAtLine(txt)
		if x1 < x {
			break
		}
		if x1 == x && st != "" {
			if st[len(st)-1:] == ":" {
				st = st[:len(st)-1]
			}
			existingSiblings[st] = true
		}
	}

	return stringPath, existingSiblings
}

//func getCurrentWordSelection(txt string, selPos int) (selStart int, selEnd int) {
//	//if selPos < len(txt) && selPos > 0 && txt[selPos:selPos+1] == "\n" {
//	//	selPos--
//	//}
//	selStart = len(txt) - 1
//	selEnd = selStart + 1
//
//	// expand left
//	for i := selPos; i > -1; i-- {
//		if i < len(txt) {
//			s := txt[i : i+1]
//			if !isLetter(s) {
//				selStart = i
//				if i != selPos {
//					selStart++
//				}
//				break
//			}
//			if i == 0 {
//				selStart--
//				break
//			}
//		}
//	}
//	// expand right
//	for i := selPos; i < len(txt); i++ {
//		s := txt[i : i+1]
//		if !isLetter(s) {
//			selEnd = i
//			break
//		}
//	}
//	if selEnd < selStart {
//		selEnd = selStart
//	}
//	return
//}

func GetLeftmostWordAtLine(r []rune) (word string, x1 int, x2 int) {
	// find the start of the word
	for i := 0; i < len(r); i++ {
		if isLetter(r[i]) {

			if r[i] == '-' {
				continue
			}

			x1 = i
			break
		}
	}
	// find the end of the word
	x2 = len(r)
	for i := x1; i < len(r); i++ {
		if !isLetter(r[i]) {
			x2 = i
			break
		}
	}
	word = string(r[x1:x2])
	return
}

// returns true for letters,numbers,UTF symbols, quotes, colons, etc.
// returns false for spaces,tabs,line feeds
func isLetter(r rune) bool {
	if r == ' ' || r == '\t' || r == '\n' {
		return false
	}
	return true
}
