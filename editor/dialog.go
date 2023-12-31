package editor

import "github.com/mrqzzz/keyboard/keys"

type Dialog interface {
	ListenKeys(key keys.Key) (stop bool, err error)
	DrawAll()
	GetTag() string
}

type DialogParent interface {
	CloseDialog(d Dialog, accept bool)
}

type Box struct {
	X      int
	Y      int
	Width  int
	Height int
}
