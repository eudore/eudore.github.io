package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var defaultconfig interface{}
var configinfos map[interface{}]*configinfo

func Reload() error {
	ReadConfig(defaultconfig)
	ReadMode(defaultconfig)
	ReadFlag(defaultconfig)
	ReadEnv(defaultconfig)
	if configinfos[defaultconfig].Test {
		indent, _ := json.MarshalIndent(defaultconfig, "", "\t")
		fmt.Println(string(indent))
		os.Exit(0)
	}
	if configinfos[defaultconfig].Help {
		Help(defaultconfig)
		os.Exit(0)
	}
	return nil
}

func SetDf(df interface{}) {
	defaultconfig = df
}

func init() {
	configinfos = make(map[interface{}]*configinfo)
	SetReload("config", 0x00, Reload)
}
