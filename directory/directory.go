package directory

import "os"

func CreateDirectoryStructure(path string) (bool, error) {
	err := os.MkdirAll(path, 0770)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateDirectory(path string) (bool, error) {
	err := os.Mkdir(path, 0770)
	if err != nil {
		return false, err
	}
	return true, nil
}
