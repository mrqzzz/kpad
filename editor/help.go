package editor

import (
	tm "github.com/buger/goterm"
	"github.com/mrqzzz/keyboard/keys"
	"strings"
)

type HelpDialog struct {
	Box
	Tag          string
	Editor       *Editor
	DialogParent DialogParent
	Help         []string
	ScrollOffset int
}

func NewHelpDialog(tag string, e *Editor, p DialogParent, x, y, width, height int) *HelpDialog {
	if x+width > e.ScreenWidth {
		x = e.ScreenWidth - width
	}
	if y+height > e.ScreenHeight-1 {
		y = e.ScreenHeight - height + 1
	}
	d := &HelpDialog{
		Tag:          tag,
		Box:          Box{x, y, width, height},
		Editor:       e,
		DialogParent: p,
	}

	st := `

F1:              This help

CTRL+space:      Kubectl autocomplete (Linux,Mac)
CTRL+k:          Kubectl autocomplete (Windows)

Arrows:          Move cursor around
Home:            Move to begin of line
End:             Move to end of line
PageUp:          Move up one page
PageDown:        Move down one page

ALT+PageUp:
CTRL+t:          Move to top of document

ALT+PageDown:
CTRL+b:          Move to bottom of document

CMD+Right:
CTRL+x:          Move to next word

CMD+Left:
CTRL+z:          Move to previous word

CTRL+d:          Delete line

CMD+Backspace:   Forward delete

CTRL+f:          Find 
CTRL+n:          Find next
`

	spl1 := strings.Split(st, "\n")
	for i, s1 := range spl1 {
		spl2 := strings.Split(s1, ":")
		spl2[0] = tm.Color(spl2[0], tm.GREEN)
		for k, _ := range spl2 {
			spl2[k] = tm.Background(spl2[k], tm.BLUE)
		}
		spl1[i] = strings.Join(spl2, "")
	}
	d.Help = spl1

	return d
}

func (d *HelpDialog) ListenKeys(key keys.Key) (stop bool, err error) {
	if key.Code == keys.Enter {
		d.DialogParent.CloseDialog(d, true)
	} else if key.Code == keys.Esc || key.Code == keys.F1 {
		d.DialogParent.CloseDialog(d, false)
	} else if key.Code == keys.PgUp {
		d.ScrollOffset = max(0, d.ScrollOffset-d.Height-2)
		d.DrawAll()
	} else if key.Code == keys.PgDown {
		d.ScrollOffset = min(len(d.Help), d.ScrollOffset+d.Height-2)
		d.DrawAll()
	} else if key.Code == keys.Up {
		d.ScrollOffset = max(0, d.ScrollOffset-1)
		d.DrawAll()
	} else if key.Code == keys.Down {
		d.ScrollOffset = min(len(d.Help), d.ScrollOffset+1)
		d.DrawAll()
	}
	return false, nil
}

func (d *HelpDialog) DrawAll() {
	rBlanks := runeRepeat(' ', d.Box.Width)
	rTitle := runeRepeat(' ', d.Box.Width)
	copy(rTitle, []rune{' ', 'H', 'e', 'l', 'p', ':'})
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
	idx := d.ScrollOffset
	for i := 0; i < d.Height-2; i++ {
		tm.MoveCursor(d.Box.X+1, d.Box.Y+i+1)
		idx++
		if idx >= len(d.Help) {
			break
		}
		tm.Print(d.Help[idx])
	}
	tm.Flush()
}

func (d *HelpDialog) GetTag() string {
	return d.Tag
}
