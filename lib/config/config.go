package config

import (
	"fmt"
	"os"
	"path/filepath"

	"weibo.com/hotweibo/util/goconfig"
)

var Instance *goconfig.ConfigFile

func LoadDir(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		Instance.AppendFiles(path)
		return nil
	})
	if err != nil {
		fmt.Printf("path err: %v\n", err)
	}
}

func init() {
	Instance = goconfig.NewConfigFile()
}
