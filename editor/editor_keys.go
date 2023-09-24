package editor

import (
	tm "github.com/buger/goterm"
	"github.com/mrqzzz/keyboard/keys"
	"strings"
)

// CTRL + Z/X or A/E : prev/next word
// CTRL + K or SPACE : kubectl explain
// CTRL + D delete row
// HOME/END PGUP/PGDOWN
// ALT + BACKSPACE : forward delete
// NOTE:
// Windows: keys.CtrlAt = Ctrl+Z
// Mac: keys.CtrlAt = Ctrl+SPACE

func (e *Editor) ListenKeys(key keys.Key) (stop bool, err error) {

	defer func() { e.LastKey = key }()

	//e.StatusBar.DrawInfo("STRING='" + key.String() + "' CODE='" + key.Code.String() + "' ALT='" + fmt.Sprint(key.AltPressed) + "' RUNES='" + fmt.Sprintf("%v", key.Runes) + "'")

	/////////////////////////////////
	// AVOID REMOTE TERMINAL JITTER PROBLEMS
	if strings.Contains(key.String(), "[") && key.AltPressed {
		return false, nil
	}
	if len(key.Runes) > 1 && areAll(key.Runes, 127) && !key.AltPressed {
		key.Runes = []rune{127}
		key.Code = keys.Backspace
	}
	if len(key.Runes) > 1 && areAll(key.Runes, 4) && !key.AltPressed {
		return false, nil
	}
	if len(key.Runes) == 1 && key.Runes[0] == 127 && key.AltPressed {
		return false, nil
	}
	/////////////////////////////////

	CTRLz := keys.CtrlZ
	if e.IsWindows {
		CTRLz = keys.CtrlAt
	}

	if e.Dialog != nil {
		return e.Dialog.ListenKeys(key)
	}

	if key.Code == keys.F1 {
		// OPEN HELP
		e.OpenHelpDialog()
	} else if key.Code == keys.CtrlC {
		return true, nil // Stop listener by returning true on Ctrl+C
	} else if key.Code == keys.Home {
		// MOVE TO BEGIN OF LINE
		e.X = 1
		//e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.End || key.Code == 91 && key.AltPressed {
		// MOVE TO END OF LINE
		//+e.X = e.ScreenWidth
		e.X = runesWidth(e.Buf[e.Y+e.Top-1])
		//e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.CtrlT || (key.Code == keys.PgUp && key.AltPressed) {
		// MOVE TO TOP OF DOCUMENT
		e.X = 1
		e.Y = 1
		e.Top = 0
		e.DrawAll()
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.CtrlB || (key.Code == keys.PgDown && key.AltPressed) {
		// MOVE TO BOTTOM OF DOCUMENT
		e.X = 1
		e.Y = 1
		e.Top = len(e.Buf) - 1
		e.DrawAll()
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.PgUp {
		// MOVE ONE PAGE UP
		e.Top -= e.ScreenHeight
		if e.Top < 0 {
			e.Top = 0
			e.Y = 1
		}
		e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		e.DrawAll()
	} else if key.Code == keys.PgDown {
		// MOVE ONE PAGE DOWN
		e.Top += e.ScreenHeight
		if e.Top > len(e.Buf)-1 {
			e.Top = len(e.Buf) - 1
		}
		e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		e.DrawAll()
	} else if key.Code == keys.Down {
		// ARROW DOWN
		if e.Y >= e.ScreenHeight && e.Top < len(e.Buf)-e.ScreenHeight {
			e.ScrollUp()
			e.StatusBar.DrawEditing()
			tm.Flush()
		}
		if e.Y < e.ScreenHeight {
			e.Y++
			e.MoveCursorSafe(e.X, e.Y)
			e.StatusBar.DrawEditing()
			tm.Flush()
		}
	} else if key.Code == keys.Up {
		// ARROW UP
		if e.Y == 1 && e.Top > 0 {
			e.ScrollDown()
			e.StatusBar.DrawEditing()
			tm.Flush()
		}
		if e.Y > 1 {
			e.Y--
			e.MoveCursorSafe(e.X, e.Y)
			e.StatusBar.DrawEditing()
			tm.Flush()
		}
	} else if key.Code == keys.Left && key.AltPressed || key.Code == keys.CtrlA || key.Code == CTRLz {
		// MOVE TO PREVIOUS WORD
		col := runesToCover(e.Buf[e.Y+e.Top-1], e.X-1)
		advances := e.GetNextWord(col, e.Y+e.Top-1, -1)
		e.CursorWithdraw(advances)
		e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.Right && key.AltPressed || key.Code == keys.CtrlE || key.Code == keys.CtrlX {
		// MOVE TO NEXT WORD
		col := runesToCover(e.Buf[e.Y+e.Top-1], e.X-1)
		advances := e.GetNextWord(col, e.Y+e.Top-1, 1)
		e.CursorAdvance(advances)
		e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.Left {
		// ARROW LEFT
		e.CursorWithdraw(-1)
		//e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.Right {
		// ARROW RIGHT
		e.CursorAdvance(1)
		//e.MoveCursorSafe(e.X, e.Y)
		e.StatusBar.DrawEditing()
		tm.Flush()
	} else if key.Code == keys.Backspace && key.AltPressed {
		// FORWARD DELETE
		row := e.Y + e.Top - 1
		col := runesToCover(e.Buf[row], e.X-1)
		if col < (len(e.Buf[row])) && e.Buf[row][col] != '\n' {
			e.DeleteAt(col+1, row)
			//e.CursorWithdraw(0)
			//e.MoveCursorSafe(e.X, e.Y)
			e.DrawAll()
		}
	} else if key.Code == keys.Backspace {
		// DELETE
		row := e.Y + e.Top - 1
		col := runesToCover(e.Buf[row], e.X-1)
		withdraws, _ := e.DeleteAt(col, row)
		e.CursorWithdraw(withdraws)
		//e.MoveCursorSafe(e.X, e.Y)
		e.DrawAll()
	} else if key.Code == keys.CtrlD {
		// DELETE ROW
		e.DeleteRow(e.Y + e.Top - 1)
		e.MoveCursorSafe(e.X, e.Y)
		e.DrawAll()
	} else if key.Code == keys.CtrlS {
		// SAVE
		e.SaveToFile()
	} else if key.Code == keys.CtrlQ {
		// QUIT
		if !e.BufferChanged || e.LastKey.Code == keys.CtrlQ {
			return true, nil
		}
		e.StatusBar.DrawError(" THERE ARE CHANGES! Press CTRL-Q again to quit anyway")
	} else if key.Code == keys.CtrlF {
		// FIND
		e.OpenSearchDialog()
	} else if key.Code == keys.CtrlN {
		// FIND NEXT
		e.FindString(e.SearchString)
	} else if key.Code == keys.CtrlK || (key.Code == keys.CtrlAt && !e.IsWindows) {
		// CALL KUBECTL
		word, _, x2 := GetLeftmostWordAtLine(e.Buf[e.Y-1+e.Top])
		if len(word) == 0 || e.X-1 <= x2 {
			e.OpenDropdown()
		}
	} else {
		// INSERT CHARACTERS
		if key.Code == keys.Enter {
			key.Runes = []rune{'\n'}
			_, x, _ := GetLeftmostWordAtLine(e.Buf[e.Y-1+e.Top])
			for i := 0; i < x; i++ {
				key.Runes = append(key.Runes, ' ')
			}
		}
		if key.Code == keys.Tab {
			key.Runes = []rune{' ', ' '}
		}

		if len(key.Runes) > 0 {

			patchKey(&key)

			// ADD TEXT
			oldY := e.Y

			row := e.Y + e.Top - 1
			col := runesToCover(e.Buf[row], e.X-1)

			insertedCharCount, rowsPushedDown := e.InsertAt(key.Runes, col, e.Y+e.Top-1)
			e.CursorAdvance(insertedCharCount)

			// optimized partial redraw:
			if e.Y == oldY {
				toIdx := e.findNextLineFeed(e.Top + e.Y - 1)
				if rowsPushedDown-1 > (toIdx - e.Top + e.Y - 1) {
					e.DrawAll()
				} else {
					e.DrawRows(e.Top+e.Y-1, min(e.Top+e.Y-1+rowsPushedDown, e.Top+e.ScreenHeight-1))
					e.StatusBar.DrawEditing()
					tm.Flush()
				}

			} else {
				e.DrawAll()
			}
		}
	}

	// Return false to continue listening
	return false, nil
}

func patchKey(key *keys.Key) {
	// some symbols are not correctly converted by the "keys" library
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
}
