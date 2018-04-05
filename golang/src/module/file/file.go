//file.go
package file;

import (
    "fmt"
    "os"
    "io"
    "strings"
    "io/ioutil"
    "net/http"
    "html/template"
    "public/config"
    "public/router"
    "public/session"
    "public/log"
)


var globalSessions *session.Manager;

const (
	Source_Local = iota
	Source_Net
	Source_Oss
	Source_Ftp
)

//  accessKeyId =   "LTAIoq1zEjIUpHUN"
//  accessKeySecret =   "CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO"
var (
    conf_accessKeyId     =   config.Getconst("file_oss_key")
    conf_accessKeySecret =   config.Getconst("file_oss_secret")
    conf_host            =   config.Getconst("file_oss_host")
    conf_upload_dir      =   config.Getconst("file_oss_upload")
)
const (
    //accessKeyId =   "LTAIoq1zEjIUpHUN"
    //accessKeySecret =   "CZ8X8rq0s7p1qjFiDba5GTIeoQJ0vO"
    //host    =   "http://wejass.oss-cn-hongkong.aliyuncs.com"
    expire_time =   60
    //upload_dir  =   "upload/"
    callbackUrl =   "http://47.52.173.119:8081/file/call"
    base64Table =   "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"  
)

func init() {
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
    mux.GetFunc("/file/:user/:zone/*",fileget)
    mux.PostFunc("/file/:user/:zone/*",fileup)
    mux.GetFunc("/file/policy",oss_policy)
    mux.PostFunc("/file/call",oss_callback)
    mux.GetFunc("/file/file",file)
}


func fileget(w http.ResponseWriter, r *http.Request) {
    log.Info(router.GetAllValues(r))
    log.Info(strings.SplitN(r.URL.Path,"/",5)[4])

    //解析模板文件
    dir := "/data/web/upload/"+strings.SplitN(r.URL.Path,"/",3)[2]
    files,err := ioutil.ReadDir(dir)
    log.Info(err)
    t, err := template.ParseFiles("/data/web/templates/file/file.html");
    log.Info(err)
    t.Execute(w, map[string]interface{}{"Files": files});
}

func fileup(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        //设置内存大小
        r.ParseMultipartForm(32 << 20); //4M
        //获取上传的文件组
        files := r.MultipartForm.File["file"];
        len := len(files);
        for i := 0; i < len; i++ {
            //打开上传文件
            file, err := files[i].Open();
            defer file.Close();
            if err != nil {
                log.Info(err);
            }
            //创建上传目录
            dir := "/data/web/upload/"+strings.SplitN(r.URL.Path,"/",3)[2]
            os.MkdirAll(dir, os.ModePerm);
            //创建上传文件
            cur, err := os.Create(dir +"/" + files[i].Filename);
            defer cur.Close();
            if err != nil {
                log.Info(err);
            }
            _, err = io.Copy(cur, file);
            if err != nil {
				fmt.Fprintf(w, "%v", "上传失败")
				return
			}
			fmt.Println("上传完成,服务器地址:",dir+files[i].Filename)
        }
    } 
}

func file(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("file"))
}
