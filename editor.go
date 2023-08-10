package main

import (
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"fmt"
	tm "github.com/buger/goterm"
	"strings"
	"time"
)

type Editor struct {
	Buf          []string
	ScreenWidth  int
	ScreenHeight int
	X            int
	Y            int
	Top          int // first visible row index
	Wordwrap     bool
}

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

	e.Buf = []string{}
	p := 0
	for i := 0; i < len(txt); i++ {
		if (i-p+1)%(e.ScreenWidth-0) == 0 || txt[i] == '\n' || i == len(txt)-1 {
			e.Buf = append(e.Buf, txt[p:i+1])
			p = i + 1
		}
	}

	e.Top = 0
	e.X = 1
	e.Y = 1

	tm.Clear() // Clear current screen

	e.DrawAll()

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {

		if key.Code == keys.CtrlC {
			return true, nil // Stop listener by returning true on Ctrl+C
		} else if key.Code == keys.Down {
			if e.Y >= e.ScreenHeight && e.Top < len(e.Buf)-e.ScreenHeight {
				// scroll up inserting the bottom line
				e.Top++
				fmt.Fprintf(tm.Screen, "\033[1S")
				tm.MoveCursor(1, e.ScreenHeight)
				st := e.Buf[e.Top+e.ScreenHeight-1]

				//if st == "\n" {
				//	st = ".\n"
				//}

				l := len(st)
				if len(st) > 0 && st[l-1] == '\n' {
					st = st[:l-1]
				}
				tm.Print(st)
				e.MoveCursorSafe(e.X, e.Y)
				tm.Flush()
				//draw()
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
				tm.MoveCursor(1, 1)
				tm.Print(e.Buf[e.Top])
				e.MoveCursorSafe(e.X, e.Y)
				tm.Flush()
				//draw()
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
		} else {
			// EDIT
			if key.Code == keys.Backspace {
				withdraws, _ := e.DeleteAt(e.X-1, e.Y+e.Top-1)
				e.CursorWithdraw(withdraws)
				e.MoveCursorSafe(e.X, e.Y)
				//e.DrawRows(e.Top+e.Y-1, e.Top+e.Y+rowsToRedraw)
				e.DrawAll()
			}
			if key.Code == keys.Enter {
				key.Runes = []rune{'\n'}
			}
			if key.Code == keys.Tab {
				key.Runes = []rune{' ', ' '}
			}

			if len(key.Runes) > 0 {
				// ADD TEXT
				oldY := e.Y
				insertedCharCount, rowsPushedDown := e.InsertAt(string(key.Runes), e.X-1, e.Y+e.Top-1)
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
				tm.Flush()

			}
		}

		return false, nil // Return false to continue listening
	})
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
			continue
		}

		st := e.Buf[n]

		// COLORIZE
		//st := strings.ReplaceAll(e.Buf[n][:w], ":", tm.Color(":", tm.BLUE))

		// full padding on all printed lines
		l := len(st)
		st = strings.Replace(st, "\n", "", -1)
		st = st + strings.Repeat(" ", e.ScreenWidth-l) + "\n"

		//// don't print the final \n on the last screen line
		l = len(st)
		if len(st) > 0 && st[l-1] == '\n' && n == e.Top+e.ScreenHeight-1 {
			st = st[:l-1]
		}

		tm.Print(st)

	}
	tm.MoveCursor(e.X, e.Y)
	tm.Flush()

}

func (e *Editor) InsertAt(ins string, col int, row int) (insertedCharCount int, rowsPushedDown int) {

	rowsPushedDown = 1
	e.ReplaceBadChars(&ins)

	if row >= len(e.Buf) {
		e.Buf = append(e.Buf, "")
	}

	var st string
	if len(e.Buf[row]) == 0 {
		st = ins
	} else {
		st = e.Buf[row][:col] + ins + e.Buf[row][col:]

	}

	p := strings.Index(st, "\n")
	if p == -1 || p == len(st)-1 {
		p = e.ScreenWidth
	} else {
		p++
		if p > e.ScreenWidth {
			p = e.ScreenWidth
		}
	}

	if len(st) < p {
		e.Buf[row] = st
	} else {
		e.Buf[row] = st[:p]
		_, rPushed := e.InsertAt(st[p:], 0, row+1)
		rowsPushedDown += rPushed
	}
	return len(ins), rowsPushedDown
}

func (e *Editor) DeleteAt(col int, row int) (numWithdraws int, rowsToRedraw int) {
	if col == 0 {
		if row > 0 {
			row2 := e.findNextLineFeed(row)
			// get the string from the cursor (at the beginning of line), down to the next \n:
			st := strings.Join(e.Buf[row:row2+1], "")
			// remove the block:
			e.Buf = append(e.Buf[:row], e.Buf[row2+1:]...)
			// insert the block at the end of the previous line (pull up)
			if len(e.Buf[row-1]) > 0 {
				e.Buf[row-1] = e.Buf[row-1][:len(e.Buf[row-1])-1]
			}
			numWithdraws = -min(len(st), e.ScreenWidth-len(e.Buf[row-1])) - 1
			_, rowsToRedraw = e.InsertAt(st, len(e.Buf[row-1]), row-1)
		}

	} else {
		// pull up
		e.Buf[row] = e.Buf[row][:col-1] + e.Buf[row][col:]
		for r := row; r < len(e.Buf); r++ {
			if len(e.Buf[r]) > 0 && e.Buf[r][len(e.Buf[r])-1:] == "\n" {
				break
			}
			if len(e.Buf)-1 > r {
				// pull up first char of next line
				if len(e.Buf[r+1]) > 0 {
					e.Buf[r] += e.Buf[r+1][:1]
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
		col++
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
		col--
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
		if strings.Index(e.Buf[i], "\n") > -1 {
			return i
		}
	}
	return len(e.Buf) - 1
}

func (e *Editor) MoveCursorSafe(x int, y int) {
	if y > len(e.Buf)-e.Top {
		y = len(e.Buf) - e.Top
	}
	st := e.Buf[e.Top+y-1]
	if x > len(st) {
		x = len(st)
	}
	if x < 0 {
		x = 0
	}
	e.X = x
	e.Y = y
	tm.MoveCursor(e.X, e.Y)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
