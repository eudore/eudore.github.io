//file.go
package file;

import (
    "fmt"
    "log"
    "os"
    "io"
    "net/http"
    "html/template"
    "public/router"
    "public/session"
)

var globalSessions *session.Manager;

const (
	Source_Local = iota
	Source_Net
	Source_Oss
	Source_Ftp
)

func init() {
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
    mux.GetFunc("/file/:user/:zone/:path",fileget)
    mux.GetFunc("/file/:user/:zone/:path/",filelist)
    //mux.PostFunc("/file/:user/:zone/:path",fileup)
    mux.GetFunc("/file/",file)
}

func filelist(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("filelist"))
}
func fileget(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("fileget"))
}
func fileup(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        //设置内存大小
        r.ParseMultipartForm(32 << 20);
        //获取上传的文件组
        files := r.MultipartForm.File["file"];
        len := len(files);
        for i := 0; i < len; i++ {
            //打开上传文件
            file, err := files[i].Open();
            defer file.Close();
            if err != nil {
                log.Fatal(err);
            }
            //创建上传目录
            os.Mkdir("./upload", os.ModePerm);
            //创建上传文件
            cur, err := os.Create("./upload/" + files[i].Filename);
            defer cur.Close();
            if err != nil {
                log.Fatal(err);
            }
            _, err = io.Copy(cur, file);
            if err != nil {
				fmt.Fprintf(w, "%v", "上传失败")
				return
			}
			fmt.Println("上传完成,服务器地址:./upload/",files[i].Filename)
        }
    } else {
        //解析模板文件
        t, _ := template.ParseFiles("./../html/test.html");
        //输出文件数据
        t.Execute(w, nil);
    }
}

func file(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("file"))
}
