// Demo code for the TextArea primitive.
package main

import (
	"mrqzzz/kpad/editor"
)

func main() {

	corporate := `123` //－23－23－23－23－23－23:
	//a－－－－－－－－－－－－
	//123`

	e := editor.NewEditor(0, 0)
	e.LoadText(corporate)
	e.Edit()

}
