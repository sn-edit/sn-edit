package file

import (
	"github.com/0x111/sn-edit/conf"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// Write The contents of the script to a file
func WriteFile(tableName string, fieldSysName string, fieldName string, extension string, contents []byte) error {
	filePath := GenerateFilePath(tableName, fieldSysName, fieldName, extension)

	if err := FileExists(filePath); err != nil {
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

// Returns error if the file exists, nil if it does not exist
func FileExists(filePath string) error {
	_, err := os.Stat(filePath)

	if !os.IsNotExist(err) {
		return err
	}

	return nil
}

func GenerateFilePath(tableName string, fieldSysName string, fieldName string, extension string) string {
	config := conf.GetConfig()
	return config.GetString("app.root_directory") + string(os.PathSeparator) + tableName + string(os.PathSeparator) + fieldSysName + string(os.PathSeparator) + fieldName + "." + extension
}