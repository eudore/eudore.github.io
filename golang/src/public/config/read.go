package config


import (
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
)


func readconfig(sour string) ([]byte,error) {
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