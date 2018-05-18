package config;
 
import (
	"fmt"
	"os"
	"strings"
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
	Session 	string
	Const       map[string]*string
	Enable      []string
	Flag 		map[string]interface{}
	Mode        map[string]interface{}
}

func (c *Config) setconf() {
	// set flag workdir 
	if value, ok := c.Flag["workdir"].(string); ok {
		conf.Workdir = value 
	}
	os.Chdir(conf.Workdir)
	// set flag config
	if value, ok := c.Flag["config"].(string); ok {
		conf.Config = value
	}	
	// get config value
	if cf, err := readconfig(c.Config);err == nil {
		err := json.Unmarshal(cf,c)
		if err != nil {
			fmt.Println("配置解析失败:", err)
		}
	}else {
		fmt.Println("配置读取失败:", err)
	}
}

func (c *Config) setmode() {	
	// set flag mode
	if value, ok := c.Flag["enable"].(string); ok {
		conf.Enable = append(conf.Enable,strings.Split(value,",")...)
		delete(c.Flag,"enable")
	}
	if value, ok := c.Flag["disable"].(string); ok {
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
}

func (c *Config) Reload() {
	c.setconf()
	c.setmode()
}

func NewConfig() *Config {
	return &Config {
		Workdir:	"/data/web",
		Tempdir:	"/data/web/template",
		Config:		"/data/web/config/conf.json",
		Pidfile:	"/var/run/index.pid",
		Logfile:	"/data/web/logs",
		IP:			"",
		Port:		8080,
		Dbconfig:	"root:@/Jass",
		Session:	`{"CookieName": "token","EnableSetCookie": true, "Gclifetime": 3600, "Maxlifetime": 3600, "Secure": true, "CookieLifeTime": 3600, "ProviderConfig": "127.0.0.1:12001"}`,
	}
}

var conf *Config

func Instance() *Config {
	return conf
}

func init() {
	// init flag
	conf = NewConfig()
	conf.Flag = readflag(os.Args[1:])
	conf.Reload()
	Reload()
	if b,err := json.Marshal(conf.Flag);err == nil {
		json.Unmarshal(b, &conf)
	}
	os.Chdir(conf.Workdir)  
}

func Reload() {
	conf.Reload()
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
