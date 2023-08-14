package editor

import (
	"atomicgo.dev/keyboard/keys"
	"fmt"
	tm "github.com/buger/goterm"
	"strconv"
)

type Dropdown struct {
	Box
	Editor        *Editor
	DialogParent  DialogParent
	Keys          []string
	Values        []string
	SelectedIndex int
	TopIndex      int
}

func NewDropdown(e *Editor, p DialogParent, x, y, width, height int, keys []string, values []string) *Dropdown {
	if x+width > e.ScreenWidth {
		x = e.ScreenWidth - width
	}
	if y+height > e.ScreenHeight {
		y = e.ScreenHeight - height
	}
	return &Dropdown{
		Box:          Box{x, y, width, height},
		DialogParent: p,
		Keys:         keys,
		Values:       values,
	}
}

func (d *Dropdown) ListenKeys(key keys.Key) (stop bool, err error) {
	if key.Code == keys.CtrlC {
		return true, nil // Stop listener by returning true on Ctrl+C
	}
	if key.Code == keys.Enter || key.Code == keys.Tab {
		d.DialogParent.CloseDialog(d, true)
	}
	if key.Code == keys.Esc {
		d.DialogParent.CloseDialog(d, false)
	}
	if key.Code == keys.Up && d.SelectedIndex > 0 {
		d.SelectedIndex--
		if d.SelectedIndex < d.TopIndex {
			d.TopIndex--
		}
		d.DrawAll()
	}
	if key.Code == keys.Down && d.SelectedIndex < len(d.Keys)-1 {
		d.SelectedIndex++
		if d.SelectedIndex >= d.TopIndex+d.Height {
			d.TopIndex++
		}
		d.DrawAll()
	}

	return false, nil
}

func (d *Dropdown) DrawAll() {
	for i := d.TopIndex; i < d.TopIndex+d.Height; i++ {
		st := fmt.Sprintf("%"+strconv.Itoa(d.Width)+"s", d.Values[i])
		if i == d.SelectedIndex {
			st = tm.Background(st, tm.BLUE)
		}
		tm.MoveCursor(d.X, d.Y+i-d.TopIndex)
		tm.Print(st)
	}
	tm.Flush()
}
