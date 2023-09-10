package main

import (
	"mrqzzz/kpad/editor"
	"os"
)

func main() {
	//
	//	corporate := `apiVersion: apps/v1
	//kind: Deployment
	//metadata:
	//  annotations:`

	e := editor.NewEditor(0, 0)
	e.Init()

	for i, arg := range os.Args {
		if i > 0 {
			if arg[0:1] != "-" {
				e.FileName = arg
			}
		}
	}

	if e.FileName != "" {
		e.LoadFromFile(e.FileName)
	}

	e.Edit()

}
