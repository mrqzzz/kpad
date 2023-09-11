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
	var err error

	e := editor.NewEditor(0, 0)
	e.InitSize()
	err = e.LoadConfig()

	for i, arg := range os.Args {
		if i > 0 {
			switch arg {
			case "-c", "--config":
				e.FileName, err = editor.GetConfigFileName()
			default:
				e.FileName = arg
			}
		}
	}

	if e.FileName != "" {
		e.LoadFromFile(e.FileName)
	}

	if err != nil {
		e.StatusBar.State = editor.StateError
		e.StatusBar.ErrorMsg = err.Error()
	}
	e.Edit()

}
