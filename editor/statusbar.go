package editor

import (
	"fmt"
	tm "github.com/buger/goterm"
	"strconv"
)

type StatusBar struct {
	Editor   *Editor
	State    State
	ErrorMsg string
	InfoMsg  string
}

func NewStatusBar(editor *Editor) *StatusBar {
	return &StatusBar{
		Editor: editor,
		State:  StateEdit,
	}
}

type State int

const (
	StateEdit State = iota
	StateFind
	StateReplace
	StateKubectl
	StateSave
	StateError
	StateInfo
)

var bufChangedChar = map[bool]string{false: " ", true: "*"}

func (s *StatusBar) Clear() {
	e := s.Editor
	tm.MoveCursor(1, e.ScreenHeight+1)
	st := fmt.Sprintf("%"+strconv.Itoa(e.ScreenWidth)+"s", "")
	tm.Print(st)
	tm.MoveCursor(e.X, e.Y)
}

func (s *StatusBar) Draw() {
	e := s.Editor
	switch s.State {
	case StateEdit:
		x := e.X
		y := e.Y + e.Top
		stCoords := fmt.Sprintf("%s%d:%d ", bufChangedChar[e.BufferChanged], x, y)
		st := fitText(e.ScreenWidth, stCoords, e.FileName)
		tm.MoveCursor(1, e.ScreenHeight+1)
		st = tm.Background(st, tm.BLUE)
		tm.Print(st)
		tm.MoveCursor(e.X, e.Y)
	case StateError:
		st := fitText(e.ScreenWidth, s.ErrorMsg, "")
		tm.MoveCursor(1, e.ScreenHeight+1)
		st = tm.Background(st, tm.RED)
		tm.Print(st)
		tm.MoveCursor(e.X, e.Y)
	case StateInfo:
		st := fitText(e.ScreenWidth, s.InfoMsg, "")
		tm.MoveCursor(1, e.ScreenHeight+1)
		st = tm.Background(st, tm.GREEN)
		tm.Print(st)
		tm.MoveCursor(e.X, e.Y)
	}
}

func (s *StatusBar) DrawEditing() {
	s.State = StateEdit
	s.Draw()
}

func (s *StatusBar) DrawInfo(infoMsg string) {
	s.InfoMsg = infoMsg
	s.State = StateInfo
	s.Draw()
	tm.Flush()
}

func (s *StatusBar) DrawError(errorMsg string) {
	s.ErrorMsg = errorMsg
	s.State = StateError
	s.Draw()
	tm.Flush()
}

func fitText(width int, leftStr string, rightStr string) string {
	remainingSpace := width - len(leftStr) - len(rightStr)
	if remainingSpace < 0 {
		trunc := len(rightStr) + remainingSpace
		if trunc < 0 {
			trunc = 0
		}
		rightStr = rightStr[:trunc]
	}

	padLen := width - len(leftStr) - len(rightStr)
	if padLen < 0 {
		padLen = 0
	}
	result := leftStr + fmt.Sprintf("%"+strconv.Itoa(padLen)+"s", "") + rightStr
	if len(result) > width {
		result = result[:width]
	}
	return result
}
