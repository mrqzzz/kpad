package editor

import (
	"os"
	"path/filepath"
)

var (
	KUBECTL = "kubectl"
)

const defaultConfig = `
#############
# flat config
#############

# "kubectl"" defines the kubectl command for completion.
# Examples:
#
#   kubectl: kubectl
#
#   kubectl: microk8s kubectl
#

kubectl: kubectl`

func GetConfigFileName() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".kpad"), nil
}

func (e *Editor) LoadConfig() error {
	fName, err := GetConfigFileName()
	if err != nil {
		return err
	}

	b, err := os.ReadFile(fName)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(fName)
			if err != nil {
				return err
			}
			_, err = f.WriteString(defaultConfig)
			if err != nil {
				return err
			}
			f.Close()

		} else {
			return err
		}
	}

	root, err := ReadFlatYaml(string(b))
	if err != nil {
		return err
	}
	for _, yNode := range root.Children {
		switch yNode.Key {
		case "kubectl":
			KUBECTL = yNode.Value
		}
	}
	return nil
}
