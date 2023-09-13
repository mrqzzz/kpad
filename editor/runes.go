package editor

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func runeReplace(r []rune, search int32, replace int32) {
	for i := 0; i < len(r); i++ {
		if r[i] == search {
			r[i] = replace
		}
	}
}

// create a rune array filled with char. if extraChar !=0, also append it, growing the array by 1
func runeRepeat(char int32, count int, extraChar int32) []rune {
	if extraChar != 0 {
		count++
	}
	res := make([]rune, count)
	for i := 0; i < count; i++ {
		res[i] = char
	}
	if extraChar != 0 {
		res[count-1] = extraChar
	}
	return res
}

func runeIndexOf(r []rune, char int32) int {
	for i := 0; i < len(r); i++ {
		if r[i] == char {
			return i
		}
	}
	return -1
}

func (e *Editor) runeReplaceBadChars(r []rune) {
	runeReplace(r, '\r', '\n')
	runeReplace(r, '\t', ' ')
}

// join the rows or runes, returning a single row of runes
func runesJoin(rows [][]rune) []rune {
	var res []rune
	for i := 0; i < len(rows); i++ {
		res = append(res, rows[i]...)
	}
	return res
}

func runeCopy(from []rune) []rune {
	res := make([]rune, len(from))
	for i := range from {
		res[i] = from[i]
	}
	return res
}

func runeCopyAppend(runes1 []rune, runes2 []rune) []rune {
	res := make([]rune, len(runes1)+len(runes2))
	idx := 0
	for i := range runes1 {
		res[idx] = runes1[i]
		idx++
	}
	for i := range runes2 {
		res[idx] = runes2[i]
		idx++
	}
	return res
}

// return the Terminal emulator Char display width of the UTF symbol
//
//	"a"  -> 1
//	"あ" -> 2
func runeWidth(r rune) int {
	if r > 255 {
		return 2
	}
	return 1
}

// return the Terminal emulator Char width of the rune array considering the Terminal emulator Char width of the UTF symbol
//
//	"abc"  -> 3
//	"あbc" -> 4
func runesWidth(r []rune) int {
	res := 0
	for i := 0; i < len(r); i++ {
		res += runeWidth(r[i])
	}
	return res
}

// return the excessing width of tha runes:
// runesExtraWidth('123')=0
// runesExtraWidth('账123')=1
// runesExtraWidth('账123账')=2
// if maxIdx>-1, ignore runes after that index
func runesExtraWidth(r []rune, maxIdx int) int {
	res := 0
	for i := 0; i < len(r); i++ {
		res += runeWidth(r[i]) - 1
		if maxIdx > -1 && i >= maxIdx {
			break
		}
	}
	return res
}

// returns how many runes of the passed array are necessary to cover the spaces, depending on their width
func runesToCover(r []rune, spaces int) int {
	w := 0
	for i := 0; i < len(r); i++ {
		w += runeWidth(r[i])
		if w > spaces {
			return i
		}
	}
	return len(r)
}

// split r in r1 and r2.  r will have at max runesWidth=spaces, r2 the remaining runes
// will also break after a \n char
func runesSplitToCover(r []rune, spaces int) (r1 []rune, r2 []rune) {
	w := 0
	for i := 0; i < len(r); i++ {
		w += runeWidth(r[i])
		if w >= spaces || (i > 0 && r[i-1] == '\n') {
			return r[:i], r[i:]
		}
	}
	return r, []rune{}
}

func areAll(r []rune, char rune) bool {
	for i := 0; i < len(r); i++ {
		if r[i] != char {
			return false
		}
	}
	return true
}
