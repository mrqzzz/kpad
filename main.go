// Demo code for the TextArea primitive.
package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/mrqzzz/tview"
	"strings"
)

func main() {

	app := tview.NewApplication()

	textArea := tview.NewTextArea().
		SetPlaceholder("Enter text here...")
	textArea.SetTitle("Text Area").SetBorder(true)
	textArea.SetWrap(false)
	helpInfo := tview.NewTextView().
		SetText(" Press F1 for help, press Ctrl-C to exit")
	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)
	pages := tview.NewPages()

	updateInfos := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfos)
	updateInfos()

	mainView := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0, 0).
		AddItem(textArea, 0, 0, 1, 2, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false).
		AddItem(position, 1, 1, 1, 1, 0, 0, false)

	help1 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Navigation

[yellow]Left arrow[white]: Move left.
[yellow]Right arrow[white]: Move right.
[yellow]Down arrow[white]: Move down.
[yellow]Up arrow[white]: Move up.
[yellow]Ctrl-A, Home[white]: Move to the beginning of the current line.
[yellow]Ctrl-E, End[white]: Move to the end of the current line.
[yellow]Ctrl-F, page down[white]: Move down by one page.
[yellow]Ctrl-B, page up[white]: Move up by one page.
[yellow]Alt-Up arrow[white]: Scroll the page up.
[yellow]Alt-Down arrow[white]: Scroll the page down.
[yellow]Alt-Left arrow[white]: Scroll the page to the left.
[yellow]Alt-Right arrow[white]:  Scroll the page to the right.
[yellow]Alt-B, Ctrl-Left arrow[white]: Move back by one word.
[yellow]Alt-F, Ctrl-Right arrow[white]: Move forward by one word.

[blue]Press Enter for more help, press Escape to return.`)
	help2 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Editing[white]

Type to enter text.
[yellow]Ctrl-H, Backspace[white]: Delete the left character.
[yellow]Ctrl-D, Delete[white]: Delete the right character.
[yellow]Ctrl-K[white]: Delete until the end of the line.
[yellow]Ctrl-W[white]: Delete the rest of the word.
[yellow]Ctrl-U[white]: Delete the current line.

[blue]Press Enter for more help, press Escape to return.`)
	help3 := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`[green]Selecting Text[white]

Move while holding Shift or drag the mouse.
Double-click to select a word.
[yellow]Ctrl-L[white] to select entire text.

[green]Clipboard

[yellow]Ctrl-Q[white]: Copy.
[yellow]Ctrl-X[white]: Cut.
[yellow]Ctrl-V[white]: Paste.
		
[green]Undo

[yellow]Ctrl-Z[white]: Undo.
[yellow]Ctrl-Y[white]: Redo.

[blue]Press Enter for more help, press Escape to return.`)
	help := tview.NewFrame(help1).
		SetBorders(4, 4, 0, 0, 2, 2)
	help.SetBorder(true).
		SetTitle("Help").
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				pages.SwitchToPage("main")
				return nil
			} else if event.Key() == tcell.KeyEnter {
				switch {
				case help.GetPrimitive() == help1:
					help.SetPrimitive(help2)
				case help.GetPrimitive() == help2:
					help.SetPrimitive(help3)
				case help.GetPrimitive() == help3:
					help.SetPrimitive(help1)
				}
				return nil
			}
			return event
		})

	drop := tview.NewDropDown()

	dropGrid := tview.NewGrid().
		AddItem(drop, 1, 1, 1, 1, 0, 0, true)

	pages.AddAndSwitchToPage("main", mainView, true)
	pages.AddPage("dropdown", dropGrid, true, false)
	pages.AddPage("help", help, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF1 {
			pages.ShowPage("help")
			return nil
		}
		return event
	})

	// CTRL+V fast paste
	//textArea.SetClipboard(nil, func() string {
	//	buf := clipboard.Read(clipboard.FmtText)
	//	return string(buf)
	//})

	dropGrid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("main")
		}
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyEnter {
			go func() {
				_, option := drop.GetCurrentOption()
				_, start, end := textArea.GetSelection()
				textArea.Replace(start, end, option)
				pages.SwitchToPage("main")
				app.Draw()
			}()
			return tcell.NewEventKey(tcell.KeyEnter, 13, tcell.ModNone)
		}
		return event
	})

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace || event.Key() == tcell.KeyTab {

			txt := textArea.GetText()
			if event.Key() == tcell.KeyCtrlSpace {
				_, selPos, _ := textArea.GetSelection()
				selStart, selEnd := getCurrentWordSelection(txt, selPos)
				textArea.Select(selStart, selEnd)

				_, x, y, _ := textArea.GetCursor()
				path := buildCurrentPath(txt, x, y)

				// populate the dropdown
				bytes, err := ExecExplain(path)
				if err != nil {
					panic(err)
				}
				root := BuildExplainFieldsTree(bytes)
				populateDropdown(drop, root)

				// position the dropdown
				rowFrom, colFrom, _, colTo := textArea.GetCursor()
				dropGrid.SetColumns(colFrom+1, colTo-colFrom, 0)
				dropGrid.SetRows(rowFrom+1, 1, 0)
				//drop.SetCurrentOption(0)

				pages.ShowPage("dropdown")
				app.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
				return nil
			} else if event.Key() == tcell.KeyTab {
				_, selStart, selEnd := textArea.GetSelection()
				textArea.Replace(selStart, selEnd, "  ")
				//app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 32, tcell.ModNone))
				//app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 32, tcell.ModNone))
				return nil
			}

		}
		return event
	})

	if err := app.SetRoot(pages,
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func populateDropdown(drop *tview.DropDown, root *Node) {
	items := []string{}
	if root != nil {
		for _, child := range root.Children {
			items = append(items, child.FieldName+"   "+child.FieldType)
		}

	}
	drop.SetOptions(items, nil)
}

func buildCurrentPath(txt string, x int, y int) string {
	var path []string
	for i := y - 1; i >= 0; i-- {
		st, x1, _ := getLeftmostWordAtLine(txt, i)
		if x1 < x && st != "" {
			if st[len(st)-1:] == ":" {
				st = st[:len(st)-1]
			}
			path = append([]string{st}, path...)
			x = x1
		}
	}
	result := strings.Join(path, ".")
	return result
}

func getCurrentWordSelection(txt string, selPos int) (selStart int, selEnd int) {
	//if selPos < len(txt) && selPos > 0 && txt[selPos:selPos+1] == "\n" {
	//	selPos--
	//}
	selStart = len(txt) - 1
	selEnd = selStart + 1

	// expand left
	for i := selPos; i > -1; i-- {
		if i < len(txt) {
			s := txt[i : i+1]
			if !isLetter(s) {
				selStart = i
				if i != selPos {
					selStart++
				}
				break
			}
			if i == 0 {
				selStart--
				break
			}
		}
	}
	// expand right
	for i := selPos; i < len(txt); i++ {
		s := txt[i : i+1]
		if !isLetter(s) {
			selEnd = i
			break
		}
	}
	if selEnd < selStart {
		selEnd = selStart
	}
	return
}

func getLeftmostWordAtLine(txt string, y int) (word string, x1 int, x2 int) {
	n := 0
	start := -1
	end := 0
	if y < 0 {
		return "", 0, 0
	}
	//  get the text at line y
	for i := 0; i < len(txt); i++ {
		if txt[i:i+1] == "\n" {
			n++
			if n == y {
				start = i
			}
			if n == y+1 {
				break
			}
		}
		end = i
	}
	st := txt[start+1 : end+1]

	// find the start of the word
	for i := 0; i < len(st); i++ {
		s := st[i : i+1]
		if isLetter(s) {

			if s == "-" {
				continue
			}

			x1 = i
			break
		}
	}
	// find the end of the word
	x2 = len(st)
	for i := x1; i < len(st); i++ {
		s := st[i : i+1]
		if !isLetter(s) {
			x2 = i
			break
		}
	}
	word = st[x1:x2]
	return
}

func isLetter(s string) bool {
	if s == " " || s == "\t" || s == "\n" {
		return false
	}
	return true
}
