package editor

import (
	"strings"
)

func populateDropdown(root *Node) {
	items := []string{}
	if root != nil {
		for _, child := range root.Children {
			items = append(items, child.FieldName+"   "+child.FieldType)
		}

	}
	//drop.SetOptions(items, nil)
}

// x and y are relative to buf, not the cursor
func BuildCurrentPath(e *Editor, x int, y int) string {
	var path []string
	for i := y; i > 0; i-- {
		fromIdx := e.findPrevLineFeed(i) + 1
		toIdx := e.findNextLineFeed(i)
		txt := runesJoin(e.Buf[fromIdx : toIdx+1])
		st, x1, _ := GetLeftmostWordAtLine(txt)
		if x1 < x && st != "" {
			if st[len(st)-1:] == ":" {
				st = st[:len(st)-1]
			}
			path = append([]string{st}, path...)
			x = x1
		}
	}
	result := strings.Join(path, ".")
	return result
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

func isLetter(r rune) bool {
	if r == ' ' || r == '\t' || r == '\n' {
		return false
	}
	return true
}
