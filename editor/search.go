package editor

import (
	"fmt"

	tm "github.com/buger/goterm"
	"github.com/mrqzzz/keyboard/keys"
)

type SearchDialog struct {
	Box
	Tag          string
	Editor       *Editor
	DialogParent DialogParent
	SearchString string
}

func NewSearchDialog(tag, searchString string, e *Editor, p DialogParent, x, y, width, height int) *SearchDialog {
	if x+width > e.ScreenWidth {
		x = e.ScreenWidth - width
	}
	if y+height > e.ScreenHeight-1 {
		y = e.ScreenHeight - height + 1
	}
	d := &SearchDialog{
		Tag:          tag,
		Box:          Box{x, y, width, height},
		Editor:       e,
		DialogParent: p,
		SearchString: searchString,
	}
	return d
}

func (d *SearchDialog) ListenKeys(key keys.Key) (stop bool, err error) {
	if key.Code == keys.CtrlC {
		return true, nil // Stop listener by returning true on Ctrl+C
	}
	if key.Code == keys.Enter {
		d.DialogParent.CloseDialog(d, true)
	} else if key.Code == keys.Esc {
		d.DialogParent.CloseDialog(d, false)
	} else if key.Code == keys.Backspace && len(d.SearchString) > 0 {
		d.SearchString = d.SearchString[:len(d.SearchString)-1]
		d.DrawAll()
	} else if len(key.Runes) > 0 {
		d.SearchString += string(key.Runes)
		d.DrawAll()
	}
	return false, nil
}

func (d *SearchDialog) DrawAll() {
	rBlanks := runeRepeat(' ', d.Box.Width)
	rTitle := runeRepeat(' ', d.Box.Width)
	copy(rTitle, []rune{'S', 'e', 'a', 'r', 'c', 'h', ':'})
	blanks := tm.Background(string(rBlanks), tm.BLUE)
	title := tm.Background(string(rTitle), tm.BLUE)
	for i := 0; i < d.Box.Height; i++ {
		tm.MoveCursor(d.X, d.Y+i)
		if i == 0 {
			tm.Print(title)
		} else {
			tm.Print(blanks)
		}
	}
	tm.MoveCursor(d.Box.X+1, d.Box.Y+1)
	tm.Print(fmt.Sprintf("%-18s", " "))
	tm.MoveCursor(d.Box.X+1, d.Box.Y+1)
	tm.Print(d.SearchString[max(0, len(d.SearchString)-18):])
	tm.Flush()
}

func (d *SearchDialog) GetTag() string {
	return d.Tag
}
