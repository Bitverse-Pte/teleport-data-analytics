package tools

import (
	"flag"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type ConfigMap struct {
	FilePath       string
	Pointer        interface{}
	LoadedCallBack func(*ConfigMap, error)
}

func InitTomlConfigs(cm []*ConfigMap) {
	if len(cm) == 0 {
		return
	}
	var defaultPath = ""
	for _, configMap := range cm {
		defaultPath += "," + configMap.FilePath
	}
	defaultPath = defaultPath[1:]
	var filePath string
	flag.StringVar(&filePath, "c", defaultPath, "init config files")
	flag.Parse()
	var paths = strings.Split(filePath, ",")
	if len(paths) != len(cm) {
		panic("-c args count error, the program needed " + strconv.Itoa(len(cm)) + " file for initial, each file split by ,")
	}
	for index, path := range paths {
		path = strings.Trim(path, " ")
		cm[index].FilePath = path
		err := DecodeToml(path, cm[index].Pointer)
		if cm[index].LoadedCallBack != nil {
			cm[index].LoadedCallBack(cm[index], err)
		}
	}
}

func DecodeToml(filepath string, pointer interface{}) error {
	_, err := toml.DecodeFile(filepath, pointer)
	return err
}
