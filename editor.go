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
	Bottom       int // last visible row index
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

	e.Draw()

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {

		if key.Code == keys.CtrlShiftUp {
			fmt.Fprintf(tm.Screen, "\033[1S")
			tm.Flush()
		} else if key.Code == keys.CtrlShiftDown {
			fmt.Fprintf(tm.Screen, "\033[1T")
			tm.Flush()
		} else if key.Code == keys.CtrlC {
			return true, nil // Stop listener by returning true on Ctrl+C
		} else if key.Code == keys.Down {
			if e.Y >= e.ScreenHeight && e.Top < len(e.Buf)-e.ScreenHeight {
				// scroll up inserting the bottom line
				e.Top++
				fmt.Fprintf(tm.Screen, "\033[1S")
				tm.MoveCursor(1, e.ScreenHeight)
				st := e.Buf[e.Top+e.ScreenHeight-1]
				if st == "\n" {
					st = ".\n"
				}

				l := len(st)
				if len(st) > 0 && st[l-1] == '\n' {
					st = st[:l-1]
				}
				tm.Print(st)
				tm.MoveCursor(e.X, e.Y)
				tm.Flush()
				//draw()
			}
			if e.Y < e.ScreenHeight {
				e.Y++
				tm.MoveCursor(e.X, e.Y)
				tm.Flush()

			}

		} else if key.Code == keys.Up {
			if e.Y == 1 && e.Top > 0 {
				// scroll down inserting the top line
				e.Top--
				fmt.Fprintf(tm.Screen, "\033[1T")
				tm.MoveCursor(1, 1)
				tm.Print(e.Buf[e.Top])
				tm.MoveCursor(e.X, e.Y)
				tm.Flush()
				//draw()
			}
			if e.Y > 1 {
				e.Y--
				tm.MoveCursor(e.X, e.Y)
				tm.Flush()
			}
		} else if key.Code == keys.Left && e.X > 1 {
			e.X--
			tm.MoveCursor(e.X, e.Y)
			tm.Flush()
		} else if key.Code == keys.Right && e.X < e.ScreenWidth {
			e.X++
			tm.MoveCursor(e.X, e.Y)
			tm.Flush()
		} else {
			// EDIT
			if key.Code == keys.Enter {
				key.Runes = []rune{'\n'}
			}
			insertedCount := e.InsertAtCursor(string(key.Runes), e.X-1, e.Y+e.Top-1)
			e.AdvanceCursor(insertedCount)
			e.Draw()
			tm.Flush()
		}

		return false, nil // Return false to continue listening
	})
	return nil
}

func (e *Editor) Draw() {
	tm.Clear()
	tm.MoveCursor(1, 1)

	for n := e.Top; n < e.Top+e.ScreenHeight; n++ {
		if n >= len(e.Buf) {
			continue
		}

		st := e.Buf[n]
		//st := strings.ReplaceAll(e.Buf[n][:w], ":", tm.Color(":", tm.BLUE))

		// don't print the final \n on the last screen line
		l := len(st)
		if len(st) > 0 && st[l-1] == '\n' && n == e.Top+e.ScreenHeight-1 {
			st = st[:l-1]
		}

		//if st == "\n" {
		//	st = ".\n"
		//}

		tm.Print(st)

	}
	tm.MoveCursor(e.X, e.Y)
	tm.Flush()

}

func (e *Editor) InsertAtCursor(ins string, col int, row int) (insertedCount int) {

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
		e.InsertAtCursor(st[p:], 0, row+1)
	}
	return len(ins)
}

func (e *Editor) AdvanceCursor(n int) {
	col := e.X - 1
	row := e.Y + e.Top - 1
	for i := 0; i < n; i++ {
		col++
		if col >= len(e.Buf[row]) {
			row++
			col = 0
		}
		if row >= len(e.Buf) {
			row = len(e.Buf)
		}
	}

	if row > e.Top+e.ScreenHeight-1 {
		e.Top = row - e.ScreenHeight + 1
		e.Draw()
	}

	e.X = col + 1
	e.Y = row - e.Top + 1
	tm.MoveCursor(e.X, e.Y)
	tm.Flush()
}

func (e *Editor) ReplaceBadChars(st *string) {
	*st = strings.ReplaceAll(*st, "\r", "\n")
	*st = strings.ReplaceAll(*st, "\t", "  ")
}
