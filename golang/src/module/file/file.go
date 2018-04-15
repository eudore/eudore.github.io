//file.go
package file;

import (
    "net/http"
    "public/config"
    "public/router"
    "public/session"
)


var globalSessions *session.Manager;
var savetype map[int]string
var savesource map[int]string
const (
	Source_Local = iota
	Source_Net
	Source_Oss
	Source_Ftp
)

//  accessKeyId =   "LTAIoq1zEjIUpHUN"
//  accessKeySecret =   "CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO"
var (
    conf_updir          =   config.Getconst("file_up_dir")
    conf_accessKeyId    =   config.Getconst("file_oss_key")
    conf_accessKeySecret=   config.Getconst("file_oss_secret")
    conf_host           =   config.Getconst("file_oss_host")
    conf_upload_dir     =   config.Getconst("file_oss_upload")
    conf_callback_domian=   config.Getconst("file_oss_callback_domain")
)

const (
    expire_time =   60
    callbackUrl =   "http://47.52.173.119:8081/file/call"
    callbackBody=   `{"filename":${object},"mimeType":${mimeType},"size":${size}}`
    callbackBodyType="application/json"
    base64Table =   "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"  
)

func init() {
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
    mux.GetFunc("/file/:user/:zone/*",local_list)
    mux.PostFunc("/file/:user/:zone/*",file_up)
    mux.PutFunc("/file/:user/:zone/*",file_up)
    mux.GetFunc("/file/policy",oss_policy)
    mux.PostFunc("/file/call",oss_callback)
    mux.GetFunc("/file/file",file)
    savetype = make(map[int]string)
    savesource = make(map[int]string)
    savetype[101] = "oss"
    savesource[101] = "http://wejass.oss-cn-hongkong.aliyuncs.com"
}


func file(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("file"))
}