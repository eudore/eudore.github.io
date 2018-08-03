package config


import (
	"os"
	"fmt"
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

// Read config
func ReadConfig(c interface{}) error {
	// change workspace
	os.Chdir(ReadSys("workdir"))
	// read config data
	ccc,_ := GetData(c,"config")
	cd, err := ReadConfigData(ccc.(string))//,ReadSys("config")
	if err != nil {
		fmt.Println("配置读取失败:", err)
		return err
	}
	// load config to c
	err = json.Unmarshal(cd,c)
	if err != nil {
		fmt.Println("配置解析失败:", err)
		return err
	}
	// save mode
	info := &configinfo{}
	json.Unmarshal(cd,info)
	configinfos[c] = info
	return nil
}

// Read mode load config
func ReadMode(c interface{}) error {
	info,ok := configinfos[c]
	if !ok {
		return nil
	}
	info.Enable = append(info.Enable, strings.Split(ReadSys("enable"), ",")...)
	info.Disable = append(info.Disable, strings.Split(ReadSys("disable"), ",")...)
	// set mode config
	for _,v := range info.getmode() {
		if b,err := json.Marshal(info.Mode[v]);err == nil && info.Mode[v] != nil {
			json.Unmarshal(b, &c)
		}
	}
	os.Chdir(info.Workdir) 
	return nil
}

// Read args to config
func ReadFlag(c interface{}) {
	var err error
	info := configinfos[c]
	for _,v := range os.Args[1:] {
		if !strings.HasPrefix(v, "--") {
			fmt.Println("invalid args",v)
			continue
		}
		kv := strings.SplitN(v[2:],"=",2)
		switch kv[0]{
		case "test","help","enable","disable","mode":
			err = SetData(info, v[2:])
		default:
			err = SetData(c, v[2:])
		}
		if err != nil {
			fmt.Println("error:",err,v)
		}
	}
}

// Read env to config
func ReadEnv(c interface{}) {
	for _, value := range os.Environ() {
		if strings.HasPrefix(value, "ENV_") {
			kv := strings.SplitN(value,"=",2)
			kv[0] = strings.ToLower(strings.Replace(kv[0],"_",".",-1))[4:]
			SetData(c,strings.Join(kv,"="))
		}
	}
}



func ReadSys(name string) string {
	// read name from env
	envname := "ENV_" + strings.ToUpper(name)
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, envname) {
			kv := append(strings.SplitN(v,"=",2),"") 
			return kv[1]
		}
	}
	// read name from flag
	flagname := "--" + name
	for _,v := range os.Args[1:] {
		if strings.HasPrefix(v, flagname) {
			kv := append(strings.SplitN(v,"=",2),"") 
			return kv[1]
		}
	}
	return ""
}

func ReadConfigData(sour string) ([]byte,error) {
	if len(sour) == 0 {
		return nil,errors.New("config is null")
	}
	s := strings.SplitN(sour, "://",2)
	switch s[0] {
	case "http":
	case "https":
		return readweb(sour)
	case "file":
		return readfile(s[1])
	default:
		return readfile(s[0])
	}
	return nil,errors.New("undefined read config: "+sour)
}

func readfile(file string) ([]byte,error) {
	return ioutil.ReadFile(file)
}

func readweb(url string) ([]byte,error) {
	resp, err := http.Get(url)
	if err!=nil {
		return nil,err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}