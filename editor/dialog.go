package editor

import "atomicgo.dev/keyboard/keys"

type Dialog interface {
	ListenKeys(key keys.Key) (stop bool, err error)
	DrawAll()
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
