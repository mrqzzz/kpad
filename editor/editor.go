package editor

import (
	"fmt"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/mrqzzz/keyboard"
	"github.com/mrqzzz/keyboard/keys"
)

type Editor struct {
	Buf          [][]rune
	ScreenWidth  int
	ScreenHeight int
	X            int
	Y            int
	Top          int // first visible row index
	Dialog       Dialog
}

var emptyDoc = []rune{'\n'}

func NewEditor(dummyX, dummyY int) *Editor {
	return &Editor{
		ScreenWidth:  dummyX,
		ScreenHeight: dummyY,
	}
}

func (e *Editor) LoadText(txt string) {
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

func (e *Editor) Init() error {
	// be sure to have a terminal
	for cnt := 0; cnt < 500; cnt++ {
		e.ScreenWidth = tm.Width()
		e.ScreenHeight = tm.Height()
		if e.ScreenWidth == 0 && e.ScreenHeight == 0 {
			time.Sleep(time.Millisecond * 10)
		} else {
			break
		}
	}

	e.Top = 0
	e.X = 1
	e.Y = 1
	return nil
}

func (e *Editor) Edit() error {
	tm.Clear() // Clear current screen
	e.DrawAll()
	keyboard.Listen(e.ListenKeys)
	return nil
}

func (e *Editor) DrawAll() {
	tm.Clear()
	e.DrawRows(e.Top, e.Top+e.ScreenHeight-1)
}

func (e *Editor) DrawRows(fromIdx int, toIdx int) {
	tm.MoveCursor(1, fromIdx-e.Top+1)
	for n := fromIdx; n <= toIdx; n++ {
		if n >= len(e.Buf) {
			break
		}

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

		st := string(runes) + "\r" // FOR WINDOWS

		// COLORIZE
		st = strings.ReplaceAll(st, ":", tm.Color(":", tm.BLUE))

		tm.Print(st)

	}
	e.MoveCursorSafe(e.X, e.Y)
	tm.Flush()
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

	return len(ins), rowsPushedDown
}

func (e *Editor) DeleteAt(col int, row int) (numWithdraws int, rowsToRedraw int) {
	if col == 0 {
		if row > 0 {
			row2 := e.findNextLineFeed(row)
			// get the string from the cursor (at the beginning of line), down to the next \n:
			st := runesJoin(e.Buf[row : row2+1])
			e.Buf = append(e.Buf[:row], e.Buf[row2+1:]...)

			extraWithdraw := 0
			// if the last rune from the previous row is a \n , then there is an extra withdraw
			if e.Buf[row-1][len(e.Buf[row-1])-1] == '\n' {
				extraWithdraw = -1
			}

			// remove last rune from the previous line
			if len(e.Buf[row-1]) > 0 {
				e.Buf[row-1] = e.Buf[row-1][:len(e.Buf[row-1])-1]
			}
			// calculate how many cursor withdraws
			emptySpaces := e.ScreenWidth - len(e.Buf[row-1]) - runesExtraWidth(e.Buf[row-1], -1)
			numWithdraws = -runesToCover(st, emptySpaces) + extraWithdraw
			//numWithdraws = -min(w1, w2) - 1
			// insert the string at the end of the previous row
			_, rowsToRedraw = e.InsertAt(st, len(e.Buf[row-1]), row-1)
		}

	} else {
		// pull up
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
		numWithdraws = -1
	}
	return
}

func (e *Editor) CursorAdvance(n int) {
	col := e.X - 1
	row := e.Y + e.Top - 1
	for i := 0; i < n; i++ {
		col++ // += runeWidth(e.Buf[row][col])
		if col >= len(e.Buf[row]) && row < len(e.Buf)-1 {
			row++
			col = 0
		}
		if row >= len(e.Buf) {
			row = len(e.Buf)
		}
	}

	if row > e.Top+e.ScreenHeight-1 {
		e.Top = row - e.ScreenHeight + 1
		e.DrawAll()
	}

	e.X = col + 1
	e.Y = row - e.Top + 1
	//tm.MoveCursor(e.X, e.Y)
	//tm.Flush()
}

func (e *Editor) CursorWithdraw(n int) {
	col := e.X - 1
	row := e.Y + e.Top - 1
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
			col = len(e.Buf[row])
		}
	}

	if row < e.Top {
		e.Top = row
		e.DrawAll()
	}
	e.X = col + 1
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
	if y > len(e.Buf)-e.Top {
		y = len(e.Buf) - e.Top
	}
	runes := e.Buf[e.Top+y-1]
	if x > len(runes) {
		x = len(runes)
	}
	if x < 0 {
		x = 0
	}
	e.X = x
	e.Y = y
	tm.MoveCursor(e.X+runesExtraWidth(runes[:e.X-1], -1), e.Y)
}

func (e *Editor) DeleteRow(idx int) {
	if len(e.Buf) == 1 {
		e.Buf[0] = runeCopy(emptyDoc)
	} else if idx < len(e.Buf) {
		e.Buf = append(e.Buf[:idx], e.Buf[idx+1:]...)
	}
}

// CTRL + Z/X or A/E : prev/next word
// CTRL + K or SPACE : kubectl explain
// CTRL + D delete row
// HOME/END PGUP/PGDOWN
// ALT + BACKSPACE : forward delete

func (e *Editor) ListenKeys(key keys.Key) (stop bool, err error) {

	if e.Dialog != nil {
		return e.Dialog.ListenKeys(key)
	}

	//tm.MoveCursor(1, e.ScreenHeight)
	//tm.Print(fmt.Sprint(key.AltPressed) + "," + string(key.Runes) + "," + strconv.Itoa(int(key.Code)) + "                  ")
	//tm.MoveCursor(e.X, e.Y)
	//tm.Flush()

	//fmt.Println(key)

	if key.Code == keys.CtrlC {
		return true, nil // Stop listener by returning true on Ctrl+C
	} else if key.Code == keys.Home {
		e.X = 1
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.CtrlK || (key.Code == keys.CtrlAt && !IsWindows()) {
		// Windows: keys.CtrlAt = Ctrl+Z
		// Mac: keys.CtrlAt = Ctrl+SPACE
		word, _, x2 := GetLeftmostWordAtLine(e.Buf[e.Y-1+e.Top])
		if len(word) == 0 || e.X-1 <= x2 {

			tm.MoveCursor(1, e.ScreenHeight)
			tm.Print("kubectl...")
			tm.MoveCursor(e.X, e.Y)
			tm.Flush()

			e.OpenDropdown()
		}

	} else if key.Code == keys.End || key.Code == 91 && key.AltPressed {
		e.X = e.ScreenWidth
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.PgUp {
		e.Top -= e.ScreenHeight
		if e.Top < 0 {
			e.Top = 0
		}
		e.MoveCursorSafe(e.X, e.Y)
		e.DrawAll()
	} else if key.Code == keys.PgDown {
		e.Top += e.ScreenHeight
		if e.Top > len(e.Buf)-1 {
			e.Top = len(e.Buf) - 1
		}
		e.MoveCursorSafe(e.X, e.Y)
		e.DrawAll()
	} else if key.Code == keys.Down {
		if e.Y >= e.ScreenHeight && e.Top < len(e.Buf)-e.ScreenHeight {
			// scroll up inserting the bottom line
			e.Top++
			fmt.Fprintf(tm.Screen, "\033[1S")
			e.MoveCursorSafe(e.X, e.Y)
			e.DrawRows(e.Top+e.ScreenHeight-1, e.Top+e.ScreenHeight-1)
		}
		if e.Y < e.ScreenHeight {
			e.Y++
			e.MoveCursorSafe(e.X, e.Y)
			tm.Flush()
		}
	} else if key.Code == keys.Up {
		if e.Y == 1 && e.Top > 0 {
			// scroll down inserting the top line
			e.Top--
			fmt.Fprintf(tm.Screen, "\033[1T")
			e.MoveCursorSafe(e.X, e.Y)
			e.DrawRows(e.Top, e.Top)
		}
		if e.Y > 1 {
			e.Y--
			e.MoveCursorSafe(e.X, e.Y)
			tm.Flush()
		}
	} else if key.Code == keys.Left && key.AltPressed || key.Code == keys.CtrlA || key.Code == keys.CtrlZ {
		advences := e.GetNextWord(e.X-1, e.Y+e.Top-1, -1)
		e.CursorWithdraw(advences)
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.Left {
		e.CursorWithdraw(-1)
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.Right && key.AltPressed || key.Code == keys.CtrlE || key.Code == keys.CtrlX {
		advences := e.GetNextWord(e.X-1, e.Y+e.Top-1, 1)
		e.CursorAdvance(advences)
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.Right {
		e.CursorAdvance(1)
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.Backspace && key.AltPressed {
		if e.Buf[e.Y+e.Top-1][e.X-1] != '\n' {
			e.DeleteAt(e.X, e.Y+e.Top-1)
			e.CursorWithdraw(0)
			e.MoveCursorSafe(e.X, e.Y)
			e.DrawAll()
		}
	} else if key.Code == keys.Backspace {
		withdraws, _ := e.DeleteAt(e.X-1, e.Y+e.Top-1)
		e.CursorWithdraw(withdraws)
		e.MoveCursorSafe(e.X, e.Y)
		e.DrawAll()
	} else if key.Code == keys.CtrlD {
		e.DeleteRow(e.Y + e.Top - 1)
		e.MoveCursorSafe(e.X, e.Y)
		e.DrawAll()
	} else {
		// EDIT
		if key.Code == keys.Enter {
			key.Runes = []rune{'\n'}
			_, x, _ := GetLeftmostWordAtLine(e.Buf[e.Y-1+e.Top])
			for i := 0; i < x; i++ {
				key.Runes = append(key.Runes, ' ')
			}

			//runes := e.Buf[e.Y-1+e.Top]
			//for i := 0; i < len(runes); i++ {
			//	if runes[i] == ' ' {
			//		key.Runes = append(key.Runes, ' ')
			//	} else {
			//		break
			//	}
			//}
		}
		if key.Code == keys.Tab {
			key.Runes = []rune{' ', ' '}
		}

		if len(key.Runes) > 0 {

			// some symbols
			if key.AltPressed && key.Runes[0] == 232 {
				key.Runes[0] = '['
			} else if key.AltPressed && key.Runes[0] == 521 {
				key.Runes[0] = ']'
			} else if key.AltPressed && key.Runes[0] == 92 {
				key.Runes[0] = '`'
			} else if key.AltPressed && key.Runes[0] == 242 {
				key.Runes[0] = '@'
			} else if key.AltPressed && key.Runes[0] == 224 {
				key.Runes[0] = '#'
			} else if key.AltPressed && key.Runes[0] == 101 {
				key.Runes[0] = '€'
			}

			// ADD TEXT
			oldY := e.Y
			insertedCharCount, rowsPushedDown := e.InsertAt(key.Runes, e.X-1, e.Y+e.Top-1)
			e.CursorAdvance(insertedCharCount)

			// optimized partial redraw:
			if e.Y == oldY {
				toIdx := e.findNextLineFeed(e.Top + e.Y - 1)
				if rowsPushedDown-1 > (toIdx - e.Top + e.Y - 1) {
					e.DrawAll()
				} else {
					e.DrawRows(e.Top+e.Y-1, min(e.Top+e.Y-1+rowsPushedDown, e.Top+e.ScreenHeight-1))
				}

			} else {
				e.DrawAll()
			}
			//e.MoveCursorSafe(e.X, e.Y)
			//tm.Flush()
		}
	}

	return false, nil // Return false to continue listening
}

func (e *Editor) OpenDropdown() {
	word, startCol, startRow, _, _ := e.GetWordAtPos(e.X-1, e.Y-1+e.Top)
	//if len(word) == 0 {
	//	word, startCol, startRow, _, _ = e.GetWordAtPos(e.X-2, e.Y-1+e.Top)
	//}
	path, existingSiblings := BuildCurrentPath(e, startCol, startRow-1)
	if path == "" {
		bytes, err := ExecKubectlApiResources()
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
		bytes, err := ExecKubectlExplain(path)
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
		}

	}
	e.Dialog = nil
	e.DrawAll()
}
