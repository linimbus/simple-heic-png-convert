package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	InputDir     string `json:"inputDir"`
	OutputDir    string `json:"outputDir"`
	PngEnable    bool   `json:"pngEnable"`
	PngCompLevel string `json:"pngCompressionLevel"`
	JpegEnable   bool   `json:"jpegEnable"`
	JpegQuality  int    `json:"jpegQuality"`
	TaskNum      int    `json:"taskNum"`
}

var configCache = Config{
	InputDir:     "",
	OutputDir:    "",
	PngEnable:    true,
	PngCompLevel: "Low",
	JpegEnable:   true,
	JpegQuality:  100,
	TaskNum:      10,
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

func InputDirSave(path string) error {
	configCache.InputDir = path
	return configSyncToFile()
}

func OutputDirSave(path string) error {
	configCache.OutputDir = path
	return configSyncToFile()
}

func PngEnableSave(flag bool) error {
	configCache.PngEnable = flag
	return configSyncToFile()
}

func PngCompressLevelSave(level string) error {
	configCache.PngCompLevel = level
	return configSyncToFile()
}

func JpegEnableSave(flag bool) error {
	configCache.JpegEnable = flag
	return configSyncToFile()
}

func JpegQualitySave(quality int) error {
	configCache.JpegQuality = quality
	return configSyncToFile()
}

func TaskNumSave(num int) error {
	configCache.TaskNum = num
	return configSyncToFile()
}

func ConfigInit() error {
	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, "config.json")

	_, err := os.Stat(configFilePath)
	if err != nil {
		err = configSyncToFile()
		if err != nil {
			logs.Error("config sync to file fail, %s", err.Error())
			return err
		}
	}

	value, err := os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		return err
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		return err
	}

	return nil
}
