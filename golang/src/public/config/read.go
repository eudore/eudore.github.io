package config


import (
	"os"
	"fmt"
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
)


func (c *Config) readflag() {
	for _,v := range os.Args[1:] {
		if !strings.HasPrefix(v, "--") {
			fmt.Println("invalid args",v)
			continue
		}
		kv := strings.SplitN(v[2:],"=",2)
		switch kv[0]{
		case "flag":
		case "mode":
			continue
		case "help":
			c.help()
			os.Exit(0)
		default:
			c.set(v[2:])
			continue
		}
		fmt.Println("error args",v)
	}
}

func (c *Config) readenv() {
	for _, value := range os.Environ() {
		if strings.HasPrefix(value, "ENV_") {
			kv := strings.SplitN(value,"=",2)
			kv[0] = strings.ToLower(strings.Replace(kv[0],"_",".",-1))[4:]
			c.set(strings.Join(kv,"="))
		}
	}
}

func readconfig(sour string) ([]byte,error) {
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