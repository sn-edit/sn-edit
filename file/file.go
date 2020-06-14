package file

import (
	"errors"
	"github.com/kennygrant/sanitize"
	log "github.com/sirupsen/logrus"
	"github.com/sn-edit/sn-edit/conf"
	"io/ioutil"
	"os"
)

// Write The contents of the script to a file
func WriteFile(tableName string, scopeName string, uniqueKeyName string, fieldName string, extension string, contents []byte) error {
	filePath := GenerateFilePath(tableName, scopeName, uniqueKeyName, fieldName, extension)

	if exists := Exists(filePath); exists == false {
		err := errors.New("file_not_found")
		log.WithFields(log.Fields{"error": err, "filepath": filePath}).Error("There was an error while checking for the file existence!")
		return err
	}

	err := ioutil.WriteFile(filePath, contents, 0644)

	if err != nil {
		log.WithFields(log.Fields{"error": err, "filepath": filePath}).Panic("There was an error while writing the file!")
		return err
	}

	return nil
}

// Read the file contents
func ReadFile(filename string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filename)
	log.WithFields(log.Fields{"file_name": filename}).Debug("Reading file contents")
	if err != nil {
		return nil, err
	}

	return dat, nil
}

// Returns error if the file exists, nil if it does not exist
func Exists(filePath string) bool {
	info, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func GenerateFilePath(tableName string, scopeName string, uniqueFieldName string, fieldName string, extension string) string {
	config := conf.GetConfig()
	return config.GetString("app.core.root_directory") + string(os.PathSeparator) + scopeName + string(os.PathSeparator) + tableName + string(os.PathSeparator) + FilterSpecialChars(uniqueFieldName) + string(os.PathSeparator) + fieldName + "." + extension
}

func FilterSpecialChars(name string) string {
	return sanitize.BaseName(name)
}
