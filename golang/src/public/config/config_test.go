package config;

import (
	"fmt"
	"os"
	"testing"
	"reflect"
	"strings"
	"public/config"
	"public/log"
)

func TestConfig(t *testing.T) {
	os.Setenv("ENV_LISTEN_HTTPS","")
	os.Setenv("ENV_LISTEN_Certfile","/data/")
	os.Setenv("ENV_LISTEN_keyfile","/data/")
	config.Reload()
	data := make(map[string]interface{})
	setMap(data,"ss.ad.dd=s")
	setMap(&data,"ss.aa=s")
	log.Json(data)
	t.Log("--end")
}

func setMap(p interface{},arg string) {
	pv := reflect.ValueOf(p)
	if pv.Kind() == reflect.Ptr {
		pv= pv.Elem()
	}
	fmt.Println(pv.Kind())
	fmt.Println(pv.Type())
	data := pv.Interface().(map[string]interface{})
	kv := append(strings.SplitN(arg,"=",2),"")
	k,v := kv[0],kv[1]
	fs := strings.Split(k,".")
	len := len(fs) - 1
	for i,n := range fs {
		if data[n] == nil {
			data[n] = make(map[string]interface{})
		}
		if i==len {
			data[n] = v
		}else{

			data = data[n].(map[string]interface{})	
		}

	}
}