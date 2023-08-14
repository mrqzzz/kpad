package editor

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"fmt"
	tm "github.com/buger/goterm"
	"strings"
	"time"
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

func (e *Editor) Edit(txt string) error {

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

	e.Top = 0
	e.X = 1
	e.Y = 1

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

		ln := len(e.Buf[n])
		extraWidth := runesExtraWidth(e.Buf[n], -1)

		var extraChar int32 = 0
		if n < e.Top+e.ScreenHeight-1 {
			extraChar = '\n'
		}

		// FULL PADDING
		runes := runeRepeat('.', max(e.ScreenWidth-extraWidth, len(e.Buf[n])), extraChar)
		copy(runes, e.Buf[n])

		if runes[ln-1] == '\n' {
			runes[ln-1] = ' '
		}

		st := string(runes)

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
			// remove the block:
			e.Buf = append(e.Buf[:row], e.Buf[row2+1:]...)
			// remove last rune from the previous line
			if len(e.Buf[row-1]) > 0 {
				e.Buf[row-1] = e.Buf[row-1][:len(e.Buf[row-1])-1]
			}
			// calculate how many cursor withdraws
			emptySpaces := e.ScreenWidth - len(e.Buf[row-1]) - runesExtraWidth(e.Buf[row-1], -1)
			numWithdraws = -runesToCover(st, emptySpaces) - 1
			numWithdraws = -runesToCover(st, emptySpaces) - 1
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

func (e *Editor) ListenKeys(key keys.Key) (stop bool, err error) {

	if e.Dialog != nil {
		return e.Dialog.ListenKeys(key)
	}

	if key.Code == keys.CtrlC {
		return true, nil // Stop listener by returning true on Ctrl+C
	} else if key.Code == keys.Home {
		e.X = 1
		e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.CtrlAt {
		e.Dialog = NewDropdown(e, e, e.X, e.Y+1, 10, 3, []string{"k1", "k2", "k3", "k4", "k5", "k6"}, []string{"ああああak1", "ak2", "ak3", "ak4", "ak5", "ak6"})
		e.Dialog.DrawAll()
	} else if key.Code == keys.End {
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
	} else if key.Code == keys.Left {
		e.CursorWithdraw(-1)
		e.MoveCursorSafe(e.X, e.Y)
		//e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.Right {
		e.CursorAdvance(1)
		e.MoveCursorSafe(e.X, e.Y)
		//e.MoveCursorSafe(e.X, e.Y)
		tm.Flush()
	} else if key.Code == keys.Backspace {
		withdraws, _ := e.DeleteAt(e.X-1, e.Y+e.Top-1)
		e.CursorWithdraw(withdraws)
		e.MoveCursorSafe(e.X, e.Y)
		//e.DrawRows(e.Top+e.Y-1, e.Top+e.Y+rowsToRedraw)
		e.DrawAll()
	} else if key.Code == keys.CtrlD {
		e.DeleteRow(e.Y + e.Top - 1)
		e.MoveCursorSafe(e.X, e.Y)
		//e.DrawRows(e.Top+e.Y-1, e.Top+e.Y+rowsToRedraw)
		e.DrawAll()
	} else {
		// EDIT
		if key.Code == keys.Enter {
			key.Runes = []rune{'\n'}
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
			e.MoveCursorSafe(e.X, e.Y)
			tm.Flush()
		}
	}

	return false, nil // Return false to continue listening
}

func (e *Editor) CloseDialog(d Dialog, accept bool) {
	if e.Dialog != nil {
		if drop, ok := e.Dialog.(*Dropdown); ok {
			if accept {
				e.InsertAt([]rune(drop.Keys[drop.SelectedIndex]), e.X-1, e.Y+e.Top-1)
			}
		}

	}
	e.Dialog = nil
	e.DrawAll()
}
