package editor

import (
	tm "github.com/buger/goterm"
	"github.com/mrqzzz/keyboard/keys"
	"strings"
)

type Dropdown struct {
	Box
	Tag           string
	Editor        *Editor
	DialogParent  DialogParent
	Keys          []string
	Values        []string
	GrayedKeys    map[string]interface{}
	SelectedIndex int
	TopIndex      int
}

func NewDropdown(tag, match string, e *Editor, p DialogParent, x, y, width, height int, keys, values []string, grayedKeys map[string]interface{}) *Dropdown {
	formatValuesStrings(values, width)
	height = min(height, len(values))
	if x+width > e.ScreenWidth {
		x = e.ScreenWidth - width
	}
	if y+height > e.ScreenHeight-1 {
		y = e.ScreenHeight - height + 1
	}
	d := &Dropdown{
		Tag:          tag,
		Box:          Box{x, y, width, height},
		Editor:       e,
		DialogParent: p,
		Keys:         keys,
		Values:       values,
		GrayedKeys:   grayedKeys,
	}
	d.selectMatch(match)
	return d
}

func formatValuesStrings(values []string, maxWidth int) {
	for i := range values {
		r := []rune(values[i])
		//w := runesWidth(r)
		//if w > maxWidth {
		r, _ = runesSplitToCover(r, maxWidth)
		//}
		w := runesWidth(r)
		if maxWidth > w {
			r = runeCopyAppend(r, runeRepeat(' ', maxWidth-w, 0))
		}
		values[i] = string(r)
		//values[i] = fmt.Sprintf("%-"+strconv.Itoa(n)+"s", string(r))
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
	for i := d.TopIndex; i < min(d.TopIndex+d.Height, len(d.Values)); i++ {
		st := d.Values[i]
		if i == d.SelectedIndex {
			st = tm.Background(st, tm.GREEN)
		} else {
			st = tm.Background(st, tm.WHITE)
		}
		if d.GrayedKeys[d.Keys[i]] != nil {
			st = tm.Color(st, tm.BLUE)
		} else {
			st = tm.Color(st, tm.BLACK)

		}

		tm.MoveCursor(d.X, d.Y+i-d.TopIndex)
		tm.Print(st)
	}
	tm.MoveCursor(d.Editor.X, d.Editor.Y)
	tm.Flush()
}

func (d *Dropdown) GetTag() string {
	return d.Tag
}

func (d *Dropdown) selectMatch(match string) {
	match = strings.ToUpper(strings.Trim(match, ":"))
	for i := range d.Values {
		if strings.Contains(strings.ToUpper(d.Values[i]), match) && match != "" {
			d.SelectedIndex = i
			if d.SelectedIndex >= d.TopIndex+d.Height {
				d.TopIndex = d.SelectedIndex - d.Height + 1
			}
			break
		}
	}
}
