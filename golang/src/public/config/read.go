package config


import (
	"fmt"
	"errors"
	"strings"
	"reflect"
	"strconv"
	"io/ioutil"
	"net/http"
)


func readflag(args []string) map[string]interface{} {
	s := reflect.TypeOf(conf).Elem()
	flag := make(map[string]interface{})
	for _,v := range args {
		if !strings.HasPrefix(v, "--") {
			fmt.Println("invalid args",v)
			continue
		}
		kv := strings.SplitN(v[2:],"=",2)
		switch kv[0]{
		case "flag":
		case "mode":
		case "disable":
			if len(kv)==2 && kv[1]!="" {
				flag[kv[0]]=kv[1]
				continue
			}
		case "help":
			for i := 0; i < s.NumField(); i++ {
				if c := s.Field(i).Tag.Get("comment");c != "" {
					fmt.Println("  --"+strings.ToLower(s.Field(i).Name),"\t",c)
				}
			}
			fmt.Println("  --help \t Show help")
			continue
		default:
			if f,ok := s.FieldByName(strings.Title(kv[0]));ok{
				if len(kv)==1 {
					kv = append(kv,"")
				}
				switch f.Type.Kind() {
				case reflect.Int:
					if i,e := strconv.Atoi(kv[1]);e==nil{
						flag[kv[0]] = i
					}else{
						fmt.Println("error args",v)
					}
				case reflect.Bool:
					if b,e := strconv.ParseBool(kv[1]);e==nil{
						flag[kv[0]] = b
					}else{
						flag[kv[0]] = true
					}
				default:
					flag[kv[0]] = kv[1]
				}
				continue
			}
		}
		fmt.Println("error args",v)
	}
	return flag
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