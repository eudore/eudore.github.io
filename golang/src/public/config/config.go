package config;
 
import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"strconv"
	"encoding/json"
)
 
type Config struct {
	Config      string      `comment:"config path"`
	Workdir     string      `comment:"Current working directory"`
	Command     string      `comment:"start command"`
	Tempdir     string      `comment:"Template file dir"`
	Pidfile     string      `comment:"Pid file path"`
	Logfile     string      `comment:"Log file path"`
	IP          string      `comment:"Listen Ip Addr"`
	Port        int         `comment:"Server use port"`
	Dbconfig    string      `comment:"MariaDB connect info"`
	Memaddr     string      `comment:"Memcached connect addr and port"`
	Const       map[string]*string
	Enable      []string
	Mode        map[string]interface{}
}

var conf *Config

func Instance() *Config {
	return conf
}

func init() {
	Reload()
}

func Reload() {
	// set default config
	conf = &Config {
		Workdir:	"/data/web",
		Tempdir:	"/data/web/template",
		Config:		"/data/web/config/conf.json",
		Pidfile:	"/var/run/index.pid",
		Logfile:	"/data/web/logs",
		IP:			"",
		Port:		8080,
		Dbconfig:	"root:@/Jass",
		Memaddr:	"127.0.0.1:11211",
	}
	// init flag
	s := reflect.TypeOf(conf).Elem()
	flag := make(map[string]interface{})
	for _,v := range os.Args[1:] {
		if !strings.HasPrefix(v, "--") {
			fmt.Println("invalid args",v)
			continue
		}
		kv := strings.SplitN(v[2:],"=",2)
		switch kv[0]{
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
	// set flag workdir 
	if value, ok := flag["workdir"].(string); ok {
		conf.Workdir = value 
	}
	os.Chdir(conf.Workdir)
	// set flag config
	if value, ok := flag["config"].(string); ok {
		conf.Config = value
	}
	// get config value
	if c,err := readconfig(conf.Config);err == nil {
		err := json.Unmarshal(c,&conf)
		if err != nil {
			fmt.Println("配置解析失败:", err)
		}
	}else {
		fmt.Println("配置读取失败:", err)
	}
	// set flag mode
	if value, ok := flag["enable"].(string); ok {
		conf.Enable = append(conf.Enable,strings.Split(value,",")...)
		delete(flag,"enable")
	}
	if value, ok := flag["disable"].(string); ok {
		var d []string
		for _,ve := range conf.Enable {
			b := true
			for _,vd := range strings.Split(value,",") {
				if ve==vd {
					b = false
					break
				}
			}
			if b{
				d=append(d,ve)
			}
		}
		conf.Enable = d
	}
	// set mode config
	for _,v := range conf.Enable {
		if b,err := json.Marshal(conf.Mode[v]);err == nil && conf.Mode[v] != nil {
			json.Unmarshal(b, &conf)
		}
	}
	// set flag config
	if b,err := json.Marshal(flag);err == nil {
		json.Unmarshal(b, &conf)
	}
	os.Chdir(conf.Workdir)  
}



func (c *Config) IsMode(m string) bool {
	for _,v := range c.Enable {
		if v==m {
			return true
		}
	}
	return false
}

func (c *Config) Getconst(k string) *string {
	if v,ok := c.Const[k]; ok {
		return v
	}
	return nil
}



func IsMode(m string) bool {
	return conf.IsMode(m)
}
	
func Getconst(k string) *string {
	return conf.Getconst(k)
}
