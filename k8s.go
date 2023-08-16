package main

import "strings"

func populateDropdown(root *Node) {
	items := []string{}
	if root != nil {
		for _, child := range root.Children {
			items = append(items, child.FieldName+"   "+child.FieldType)
		}

	}
	//drop.SetOptions(items, nil)
}

func buildCurrentPath(txt string, x int, y int) string {
	var path []string
	for i := y - 1; i >= 0; i-- {
		st, x1, _ := getLeftmostWordAtLine(txt, i)
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

func getCurrentWordSelection(txt string, selPos int) (selStart int, selEnd int) {
	//if selPos < len(txt) && selPos > 0 && txt[selPos:selPos+1] == "\n" {
	//	selPos--
	//}
	selStart = len(txt) - 1
	selEnd = selStart + 1

	// expand left
	for i := selPos; i > -1; i-- {
		if i < len(txt) {
			s := txt[i : i+1]
			if !isLetter(s) {
				selStart = i
				if i != selPos {
					selStart++
				}
				break
			}
			if i == 0 {
				selStart--
				break
			}
		}
	}
	// expand right
	for i := selPos; i < len(txt); i++ {
		s := txt[i : i+1]
		if !isLetter(s) {
			selEnd = i
			break
		}
	}
	if selEnd < selStart {
		selEnd = selStart
	}
	return
}

func getLeftmostWordAtLine(txt string, y int) (word string, x1 int, x2 int) {
	n := 0
	start := -1
	end := 0
	if y < 0 {
		return "", 0, 0
	}
	//  get the text at line y
	for i := 0; i < len(txt); i++ {
		if txt[i:i+1] == "\n" {
			n++
			if n == y {
				start = i
			}
			if n == y+1 {
				break
			}
		}
		end = i
	}
	st := txt[start+1 : end+1]

	// find the start of the word
	for i := 0; i < len(st); i++ {
		s := st[i : i+1]
		if isLetter(s) {

			if s == "-" {
				continue
			}

			x1 = i
			break
		}
	}
	// find the end of the word
	x2 = len(st)
	for i := x1; i < len(st); i++ {
		s := st[i : i+1]
		if !isLetter(s) {
			x2 = i
			break
		}
	}
	word = st[x1:x2]
	return
}

func isLetter(s string) bool {
	if s == " " || s == "\t" || s == "\n" {
		return false
	}
	return true
}
