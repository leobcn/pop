package pop

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Corn map[string]interface{}

// Generate populates a new tree of files on disk populated with the given
// `Corn` (map). The root directory path is returned as a string. If an error
// occured during the generation, a non-nil error is returned.
func Generate(files Corn) (root string, err error) {
	root, err = ioutil.TempDir(root, "pop")
	if err != nil {
		err = fmt.Errorf("pop: cannot generate root directory: %s", err)
		return
	}

	err = GenerateFromRoot(root, files)
	return
}

// GenerateFromRoot populates the given `root` directory with the given `Corn`
// (map). If an error occured during the generation, a non-nil error is
// returned.
func GenerateFromRoot(root string, files Corn) (err error) {
	if root == "" {
		return fmt.Errorf("pop: root directory cannot be nil", err)
	}

	if err = createDir(root); err != nil {
		return err
	}

	for name, content := range files {
		if err = generate(root, name, content); err != nil {
			return err
		}
	}

	return nil
}

func createDir(path string) error {
	if err := os.MkdirAll(path, 0700); err != nil {
		return fmt.Errorf("pop: cannot create directory %s: %s", path, err)
	}

	return nil
}

func generate(root string, name string, content interface{}) error {
	if strings.HasSuffix(name, "/") {
		return generateDir(root, name, content)
	} else {
		return generateFile(root, name, content)
	}
}

func generateDir(root string, name string, content interface{}) error {
	var err error
	dirPath := path.Join(root, name)

	// Create the directory
	if err = createDir(dirPath); err != nil {
		return err
	}

	// Return without error if the content is nil
	if content == nil {
		return nil
	}

	// Generate the directory content
	switch content := content.(type) {
	case Corn:
		for subName, subContent := range content {
			if err = generate(dirPath, subName, subContent); err != nil {
				return err
			}
		}
		return nil

	default:
		return fmt.Errorf("pop: directory content is typed %T instead of Corn", content)
	}
}

func generateFile(root string, name string, content interface{}) error {
	filePath := path.Join(root, name)

	// Open the file, which must not already exist
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("pop: cannot create file %s: %s", filePath, err)
	}
	defer f.Close()

	// Generate the content only if `Content` is non-nil or a non-empty string
	if content != nil {
		text, ok := content.(string)
		if !ok {
			return fmt.Errorf("pop: file content of %s should be a string", filePath)
		}

		if _, err = f.WriteString(text); err != nil {
			return fmt.Errorf("pop: cannot write file %s: %s", filePath, err)
		}
	}

	return nil
}