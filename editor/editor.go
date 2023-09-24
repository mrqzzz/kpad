package editor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/mrqzzz/keyboard"
	"github.com/mrqzzz/keyboard/keys"
)

type Editor struct {
	FileName      string
	Buf           [][]rune
	ScreenWidth   int
	ScreenHeight  int
	X             int
	Y             int
	Top           int // first visible row index
	Dialog        Dialog
	StatusBar     *StatusBar
	BufferChanged bool
	LastKey       keys.Key
	IsWindows     bool
	SearchString  string
}

var emptyDoc = []rune{'\n'}

func NewEditor(dummyX, dummyY int) *Editor {
	e := &Editor{
		ScreenWidth:  dummyX,
		ScreenHeight: dummyY,
	}
	e.StatusBar = NewStatusBar(e)
	e.Buf = append(e.Buf, runeCopy(emptyDoc))
	e.IsWindows = IsWindows()
	return e
}

func (e *Editor) StringToBuf(txt string) {
	// prepare buffer
	txtRune := []rune(txt)
	e.Buf = [][]rune{}
	w := 0
	p := 0
	l := len(txtRune)
	for i := 0; i < l; i++ {
		w += runeWidth(txtRune[i])
		if (w+1) > e.ScreenWidth-1 || txtRune[i] == '\n' || i == l-1 {
			e.Buf = append(e.Buf, txtRune[p:i+1])
			p = i + 1
			w = 0
		}
	}
	// empty doc?
	if len(e.Buf) == 0 {
		e.Buf = append(e.Buf, runeCopy(emptyDoc))
	}
	// be sure the last row has a LF at the end
	lastRow := e.Buf[len(e.Buf)-1]
	if lastRow[len(lastRow)-1] != '\n' {
		e.Buf[len(e.Buf)-1] = runeCopyAppend(e.Buf[len(e.Buf)-1], []rune{'\n'})
	}
}

func (e *Editor) BufToString() string {
	buf := bytes.Buffer{}
	for _, runes := range e.Buf {
		buf.WriteString(string(runes))
	}
	return buf.String()
}

func (e *Editor) LoadFromFile(fileName string) {
	e.FileName = fileName
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		e.StatusBar.DrawError(err.Error())
		return
	}
	e.StringToBuf(string(bytes))
}

func (e *Editor) SaveToFile() {

	f, err := os.Create(e.FileName)
	if err != nil {
		e.StatusBar.DrawError(err.Error())
		return
	}
	defer f.Close()

	//var row  []rune{}
	for _, runes := range e.Buf {
		_, err = f.WriteString(string(runes))
	}
	e.BufferChanged = false
	e.StatusBar.DrawInfo("Saved " + e.FileName)
}

func (e *Editor) InitSize() error {
	// be sure to have a terminal
	maxAttempts := 500

	for cnt := 1; cnt <= maxAttempts; cnt++ {
		e.ScreenWidth = tm.Width()
		e.ScreenHeight = tm.Height() - 1
		if e.ScreenWidth <= 1 || e.ScreenHeight <= 1 {
			time.Sleep(time.Millisecond * 10)
		} else {
			break
		}
		if cnt == maxAttempts {
			return errors.New("cannot initialise terminal")
		}
	}

	e.Top = 0
	e.X = 1
	e.Y = 1

	return nil
}

func (e *Editor) DetectSizeChange() (quit chan interface{}) {
	quit = make(chan interface{})
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-time.After(time.Second):
				w := tm.Width()
				h := tm.Height() - 1
				if w != e.ScreenWidth || h != e.ScreenHeight {
					e.ScreenWidth = w
					e.ScreenHeight = h
					e.StringToBuf(e.BufToString())
					e.MoveCursorSafe(e.X, e.Y)
					e.DrawAll()
				}
			}

		}
	}()
	return quit
}

func (e *Editor) Edit() {
	e.DrawAll()
	quitChan := e.DetectSizeChange()
	keyboard.Listen(e.ListenKeys)
	quitChan <- 0
	tm.Clear()
	tm.Flush()
}

func (e *Editor) DrawAll() {
	//tm.Clear()
	e.DrawRows(e.Top, e.Top+e.ScreenHeight-1)
	e.StatusBar.Draw()
	tm.Flush()
}

func (e *Editor) DrawRows(fromIdx int, toIdx int) {
	for n := fromIdx; n <= toIdx; n++ {
		tm.MoveCursor(1, n-e.Top+1)
		var st string
		if n < len(e.Buf) {
			runes := runeCopy(e.Buf[n])

			l := len(runes)
			if n < toIdx && l+runesExtraWidth(runes, -1) >= e.ScreenWidth-1 && runes[l-1] != '\n' {
				runes = append(runes, '\n')
			}

			if n == toIdx {
				l := len(runes)
				if runes[l-1] == '\n' {
					runes = runes[0 : l-1]
				}
			}
			// COLORIZE
			st = e.colorize(runes, n)
		} else {
			st = "\n"
		}

		var st1 string
		var st2 string
		i := strings.IndexByte(st, '\n')
		if i > -1 {
			st1 = st[:i]
			st2 = st[i:]
		} else {
			st1 = st
			st2 = ""
		}

		// string + clear to EOL + \r (\r is for windows)
		out := st1 + "\033[0K" + st2 // + "\r"
		tm.Print(out)
	}
	//e.MoveCursorSafe(e.X, e.Y)
}

func (e *Editor) colorize(r []rune, row int) string {
	if len(r) == 0 {
		return ""
	}
	st := string(r)

	// get the full wrapped line (from previous \n to current row)
	fromIdx := e.findPrevLineFeed(row) + 1
	toIdx := e.findNextLineFeed(row)
	wrappedLine := runesJoin(e.Buf[fromIdx : toIdx+1])
	word, _, _ := GetLeftmostWordAtLine(wrappedLine)
	if len(word) > 0 && word[0] == '#' {
		// COMMENT (also in wrapped lines)
		l := len(st)
		if st[l-1] == '\n' {
			l--
		}
		st = tm.HighlightRegion(st, 0, l, tm.CYAN)
	} else {
		// fields at the beginning of the line. e.g.: "metadata:"
		_, x1, x2 := GetLeftmostWordAtLine(r)
		if x2 > x1 && r[x2-1] == ':' {
			isField := true
			for i := x1; i < x2-1; i++ {
				if !isAlphanumeric(r[i]) {
					isField = false
					break
				}
			}
			if isField {
				st = tm.HighlightRegion(st, x1, x2-1, tm.BLUE)
			}
		}
		// all colons:
		st = strings.ReplaceAll(st, ":", tm.Color(":", tm.MAGENTA))

		// curly braces
		st = strings.ReplaceAll(st, "{", tm.Color("{", tm.YELLOW))
		st = strings.ReplaceAll(st, "}", tm.Color("}", tm.YELLOW))
	}

	return st
}

func (e *Editor) InsertAt(ins []rune, col int, row int) (insertedCharCount int, rowsPushedDown int) {

	if len(ins) == 0 {
		return 0, 0
	}

	rowsPushedDown = 1
	e.runeReplaceBadChars(ins)
	//e.ReplaceBadChars(&ins)

	if row >= len(e.Buf) {
		e.Buf = append(e.Buf, []rune{})
	}

	var st []rune
	if len(e.Buf[row]) == 0 {
		st = runeCopy(ins)
	} else {
		st = runeCopyAppend(e.Buf[row][:col], ins)
		st = runeCopyAppend(st, e.Buf[row][col:])

	}

	st1, st2 := runesSplitToCover(st, e.ScreenWidth)
	e.Buf[row] = runeCopy(st1)
	_, rPushed := e.InsertAt(runeCopy(st2), 0, row+1)
	rowsPushedDown += rPushed

	e.BufferChanged = true

	return len(ins), rowsPushedDown
}

func (e *Editor) DeleteAt(col int, row int) (numWithdraws int, rowsToRedraw int) {
	if col == 0 {
		if row > 0 {
			row2 := e.findNextLineFeed(row)
			// get the string from the cursor (at the beginning of line), down to the next \n:
			block := runesJoin(e.Buf[row : row2+1])
			// delete the line from Buf
			e.Buf = append(e.Buf[:row], e.Buf[row2+1:]...)

			// remove last rune from the previous line
			if len(e.Buf[row-1]) > 0 {
				e.Buf[row-1] = e.Buf[row-1][:len(e.Buf[row-1])-1]
			}
			// calculate how many cursor withdraws
			emptySpaces := e.ScreenWidth - runesWidth(e.Buf[row-1])
			numWithdraws = -runesToCover(block, emptySpaces-1)
			//numWithdraws = -min(w1, w2) - 1
			// insert the string at the end of the previous row
			_, rowsToRedraw = e.InsertAt(block, len(e.Buf[row-1]), row-1)
			if numWithdraws != 0 {
				e.BufferChanged = true
			}
		}

	} else {
		// pull up

		// withdraws logic depending on char widths
		w1 := runeWidth(e.Buf[row][col-1])
		w2 := runeWidth(e.Buf[row][col])
		numWithdraws = -w1
		if w1 == 1 && w2 == 2 {
			numWithdraws = 0
		} else if w1 == 2 && w2 == 2 {
			numWithdraws = -1
		} else if w1 == 2 && e.Buf[row][col] == '\n' {
			numWithdraws = -1
		}

		e.Buf[row] = append(e.Buf[row][:col-1], e.Buf[row][col:]...)
		for r := row; r < len(e.Buf); r++ {
			if len(e.Buf[r]) > 0 && e.Buf[r][len(e.Buf[r])-1] == '\n' {
				break
			}
			if len(e.Buf)-1 > r {
				// pull up first char of next line
				if len(e.Buf[r+1]) > 0 {
					e.Buf[r] = append(e.Buf[r], e.Buf[r+1][0])
					e.Buf[r+1] = e.Buf[r+1][1:]
					if len(e.Buf[r+1]) == 0 {
						// remove this row
						e.Buf = append(e.Buf[:r+1], e.Buf[r+2:]...)
						break

					}
				}
			}
			rowsToRedraw++
		}
		e.BufferChanged = true
	}
	return
}

func (e *Editor) CursorAdvance(n int) {
	row := e.Y + e.Top - 1
	col := runesToCover(e.Buf[row], e.X-1)

	for i := 0; i < n; i++ {
		col++
		if col >= len(e.Buf[row]) {
			if row >= len(e.Buf)-1 {
				col--
			} else {
				row++
				col = 0
			}
		}
		if row >= len(e.Buf) {
			row = len(e.Buf) - 1
		}
	}

	if row > e.Top+e.ScreenHeight-1 {
		e.Top = row - e.ScreenHeight + 1
		e.DrawAll()
	}

	e.X = runesWidth(e.Buf[row][:col]) + 1 // col + 1
	e.Y = row - e.Top + 1
	//tm.MoveCursor(e.X, e.Y)
	//tm.Flush()
}

func (e *Editor) CursorWithdraw(n int) {
	row := e.Y + e.Top - 1
	col := 0
	if e.X > 1 {
		col = runesToCover(e.Buf[row], e.X-1)
	}
	for i := 0; i > n; i-- {
		col-- // -= runeWidth(e.Buf[row][col])
		if col < 0 && row == 0 {
			col = 0
		}
		if col < 0 {
			row--
			if row < 0 {
				row = 0
			}
			col = len(e.Buf[row]) - 1
		}
	}

	if row < e.Top {
		e.Top = row
		e.DrawAll()
	}
	//e.X = col + 1

	e.X = runesWidth(e.Buf[row][:col]) + 1
	e.Y = row - e.Top + 1
	//tm.MoveCursor(e.X, e.Y)
	//tm.Flush()
}

func (e *Editor) ReplaceBadChars(st *string) {
	*st = strings.ReplaceAll(*st, "\r", "\n")
	*st = strings.ReplaceAll(*st, "\t", "  ")
}

// returns the first row index in Buf where a \n is found, starting from fromIndex, included.
func (e *Editor) findNextLineFeed(fromIdx int) int {
	for i := fromIdx; i < len(e.Buf); i++ {
		if runeIndexOf(e.Buf[i], '\n') > -1 {
			return i
		}
	}
	return len(e.Buf) - 1
}

// returns the row index in Buf where a \n is found, going back, starting from fromIndex, excluded.
func (e *Editor) findPrevLineFeed(fromIdx int) int {
	for i := fromIdx - 1; i >= 0; i-- {
		if runeIndexOf(e.Buf[i], '\n') > -1 {
			return i
		}
	}
	return -1
}

func (e *Editor) GetWordAtPos(col, row int) (word []rune, startCol, startRow, endCol, endRow int) {
	if col < 0 {
		col = 0
	}
	if row > len(e.Buf)-1 {
		return
	}
	startCol = col
	startRow = row
	c := col
	r := row
	//find left limit
	for {
		if !isLetter(e.Buf[r][c]) && c != col {
			break
		}
		startCol = c
		startRow = r
		c--
		if c < 0 {
			r--
			if r < 0 {
				break
			}
			c = len(e.Buf[r]) - 1
		}
	}
	endCol = col - 1
	endRow = row
	c = col
	r = row
	//find right limit
	for {
		if !isLetter(e.Buf[r][c]) {
			break
		}
		endCol = c
		endRow = r
		c++
		if c >= len(e.Buf[r]) {
			r++
			if r >= len(e.Buf) {
				break
			}
			c = 0
		}
	}
	if startRow != endRow {
		word = runeCopyAppend(e.Buf[startRow][startCol:], e.Buf[endRow][:endCol+1])
	} else {
		word = runeCopy(e.Buf[startRow][startCol : endCol+1])
	}
	return
}

func (e *Editor) GetNextWord(col, row int, increment int) (advancements int) {
	insideWord := true
	if col >= len(e.Buf[row]) {
		col = len(e.Buf[row]) - 1
	}
	for {
		onLetter := isLetter(e.Buf[row][col])
		if !insideWord && onLetter {
			return
		}
		if !onLetter {
			insideWord = false
		}
		col += increment
		if col >= len(e.Buf[row]) {
			col = 0
			row += increment
			if row >= len(e.Buf) || row < 0 {
				return
			}
		}
		if col < 0 {
			row += increment
			if row < 0 {
				return
			}
			col = len(e.Buf[row]) - 1
		}
		advancements += increment
	}
}

func (e *Editor) MoveCursorSafe(x int, y int) {

	if e.Top > len(e.Buf) {
		e.Top = len(e.Buf) - 1
	}
	if y < 1 {
		y = 1
		e.Top--
	}

	if y > len(e.Buf)-e.Top {
		y = len(e.Buf) - e.Top
	}
	runes := e.Buf[e.Top+y-1]
	spaces := runesWidth(runes)
	if x > spaces {
		x = spaces
	}

	// properly select a wide char (not half)
	w := 0
	for i := 0; i <= x; i++ {
		if i >= len(runes) {
			x = w
			break
		}
		rw := runeWidth(runes[i])
		w += rw
		if w > x {
			x = w - 1
			if i > 0 && rw == 1 && runeWidth(runes[i-1]) == 2 {
				x++
			}
			break
		}
	}

	if x < 1 {
		x = 1
	}

	e.X = x
	e.Y = y
}

func (e *Editor) DeleteRow(idx int) {
	e.BufferChanged = true
	if len(e.Buf) == 1 {
		e.Buf[0] = runeCopy(emptyDoc)
	} else if idx < len(e.Buf) {
		e.Buf = append(e.Buf[:idx], e.Buf[idx+1:]...)
	}
}

func (e *Editor) OpenHelpDialog() {
	e.Dialog = NewHelpDialog("help", e, e, 2, 2, max(1, e.ScreenWidth-2), max(1, e.ScreenHeight-3))
	e.Dialog.DrawAll()
}

func (e *Editor) OpenDropdown() {
	word, startCol, startRow, _, _ := e.GetWordAtPos(e.X-1, e.Y-1+e.Top)
	path, existingSiblings := BuildCurrentPath(e, startCol, startRow-1)
	if path == "" {
		bytes, err := e.ExecKubectlApiResources()
		if err == nil {
			resourceNames, apiVersions := BuildApiResourcesList(bytes)
			keys := []string{}
			values := []string{}
			for i := range resourceNames {
				val := fmt.Sprintf("%.25s", fmt.Sprintf("%-25s", resourceNames[i])) + " " + fmt.Sprintf("%.15s", fmt.Sprintf("%-15s", apiVersions[i]))
				keys = append(keys, resourceNames[i]+":"+apiVersions[i])
				values = append(values, val)
			}
			e.Dialog = NewDropdown("api-resources", string(word), e, e, e.X, e.Y+1, min(40, e.ScreenWidth), min(16, e.ScreenHeight), keys, values, existingSiblings)
			e.Dialog.DrawAll()
		}
	} else {
		// there is a path like "pod.metadata"
		bytes, err := e.ExecKubectlExplain(path)
		if err == nil {
			root := BuildExplainFieldsTree(bytes)
			if root != nil {
				keys := []string{}
				values := []string{}
				for _, child := range root.Children {
					val := fmt.Sprintf("%.25s", fmt.Sprintf("%-25s", child.FieldName)) + " " + fmt.Sprintf("%.15s", fmt.Sprintf("%-15s", child.FieldType))
					keys = append(keys, child.FieldName)
					values = append(values, val)
				}
				e.Dialog = NewDropdown("explain", string(word), e, e, e.X, e.Y+1, min(40, e.ScreenWidth), min(16, e.ScreenHeight), keys, values, existingSiblings)
				e.Dialog.DrawAll()
			}
		}
	}
}

func (e *Editor) OpenSearchDialog() {
	e.Dialog = NewSearchDialog("search", e.SearchString, e, e, 2, 2, min(20, e.ScreenWidth-4), 3)
	e.Dialog.DrawAll()
}

func (e *Editor) CloseDialog(d Dialog, accept bool) {
	if e.Dialog != nil {
		if drop, ok := e.Dialog.(*Dropdown); ok {
			if accept {
				word, startCol, _, col, row := e.GetWordAtPos(e.X-1, e.Y-1+e.Top)

				// if there is a space after the word, delete it because it will be added inserting the completion
				runes := e.Buf[e.Y-1+e.Top]
				if len(runes) > col+1 && runes[col+1] == ' ' {
					col++
					word = append(word, ' ')
				}

				// delete each rune
				delta := 0
				if len(word) > 0 {
					delta = e.X - 1 - startCol
				}
				for i := 0; i < len(word); i++ {
					e.DeleteAt(col+1, row)
					col--
					if col < 0 {
						break
					}
				}
				e.CursorWithdraw(-delta)
				e.MoveCursorSafe(e.X, e.Y)
				switch d.GetTag() {
				case "api-resources":
					st := strings.Split(drop.Keys[drop.SelectedIndex], ":")
					template := []rune(GenerateResourceTemplate(st[0], st[1]))
					e.InsertAt(template, e.X-1, e.Y+e.Top-1)
					e.CursorAdvance(len(template))
				case "explain":
					st := drop.Keys[drop.SelectedIndex] + ": "
					e.InsertAt([]rune(st), e.X-1, e.Y+e.Top-1)
					e.CursorAdvance(len(st))

				}
			}
			e.Dialog = nil
			e.DrawAll()
		} else if searchDialog, ok := e.Dialog.(*SearchDialog); ok {
			e.Dialog = nil
			if accept {
				e.SearchString = searchDialog.SearchString
				e.FindString(e.SearchString)
			} else {
				e.DrawAll()
			}
		} else if _, ok := e.Dialog.(*HelpDialog); ok {
			e.Dialog = nil
			e.DrawAll()
		}
	}
}

func (e *Editor) FindString(searchString string) {
	cnt := 0
	var p int
	var st string
	searchString = strings.ToLower(searchString)
	for i := e.Top + e.Y - 1; i < len(e.Buf); i++ {
		// get the whole line from the current row down to the next \n
		n := e.findNextLineFeed(i) + 1
		runes := runesJoin(e.Buf[i:n])
		if cnt == 0 {
			st = string(runes[e.X:])
		} else {
			st = string(runes)
		}
		p = strings.Index(strings.ToLower(st), searchString)
		offs := iifInt(cnt == 0, e.X, 0)
		if p > -1 && p < e.ScreenWidth-offs-1 {
			if cnt == 0 {
				e.X = e.X + p + 1
			} else {
				e.X = p + 1
			}
			if i >= e.ScreenHeight+e.Top {
				e.Top = i - e.ScreenHeight + 1
				e.Y = e.ScreenHeight
			} else {
				e.Y = e.Y + cnt
			}
			e.MoveCursorSafe(e.X, e.Y)
			e.DrawAll()
			return
		}
		cnt++
	}
	e.DrawAll()
	e.StatusBar.DrawError(fmt.Sprintf(`"%s" Not found`, e.SearchString))
	return
}

func (e *Editor) ScrollUp() {
	e.Top++

	// scroll up excluding the top row to prevent the backscroll buffer to fill
	// set the scroll region        : "\033[<top row>;<bottom row>r"
	// scroll the region up 1 line  : "\033[1S"
	// reset the scroll region      : "\033[r"
	tm.Print("\033[2;"+strconv.Itoa(e.ScreenHeight)+"r", "\033[1S", "\033[r")

	// draw the top row
	tm.MoveCursor(1, 1)
	e.DrawRows(e.Top, e.Top)

	// draw the new last row
	e.MoveCursorSafe(e.X, e.Y)
	e.DrawRows(e.Top+e.ScreenHeight-1, e.Top+e.ScreenHeight-1)

}

func (e *Editor) ScrollDown() {
	// scroll down inserting the top line
	e.Top--

	// scroll down
	tm.Print("\033[1T")

	// draw the new top row
	e.MoveCursorSafe(e.X, e.Y)
	e.DrawRows(e.Top, e.Top)

}
