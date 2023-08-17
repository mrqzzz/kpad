// Demo code for the TextArea primitive.
package main

import (
	"mrqzzz/kpad/editor"
)

func main() {

	corporate := `apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:`

	e := editor.NewEditor(0, 0)
	e.Init()
	e.LoadText(corporate)
	e.Edit()

}
