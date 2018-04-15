package file;

import (
    "fmt"
    "os"
    "io"
    "path"
    "strings"
    "net/http"
    "encoding/json"
    "public/router"
    "public/log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)
//type uphandle func(http.ResponseWriter, *http.Request)

var uptype map[string]func(http.ResponseWriter, *http.Request)

func init() {
    uptype = make(map[string]func(http.ResponseWriter, *http.Request) )
    uptype["localmore"]     =   up_localmore
    uptype["localone"]      =   up_localone
    uptype["localmulti"]    =   up_localmulti
    uptype["localpart"]     =   up_localmulti
    uptype["ossone"]        =   up_ossone
    uptype["ossmulti"]      =   up_localmulti
    uptype["osscallback"]   =   up_localmulti
}

func file_up(w http.ResponseWriter, r *http.Request) {
    if r.URL.RawQuery != "" {  
        if v,ok := uptype[r.URL.RawQuery]; ok {
            v(w,r)
        }else{
            log.Warning("invalid upload file type: ",r.URL.RawQuery)
        }
    }else{
        file_selectsave(w,r)
    }
}

func file_selectsave(w http.ResponseWriter, r *http.Request) {
    if db, err := sql.Open("mysql","root:@/Jass");err==nil {
        defer db.Close()
        var s int;
        args := router.GetAllValues(r)
        err := db.QueryRow("SELECT Source FROM v_auth_authorized WHERE UName=? and PName=?;",args["user"],args["zone"]).Scan(&s)
        log.Info(err)
        if err==nil {
            responseBody,_ := json.Marshal(map[string]interface{}{"result": 0,"type": savetype[s],"data": s})
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            w.Write(responseBody)
            return
        }
        log.Info(args["user"],args["zone"],s)
    }       
}

func file_insert(path ,source string) error {
    if db, err := sql.Open("mysql","root:@/Jass");err==nil {
        defer db.Close()
        if stmt, err := db.Prepare("INSERT tb_note_save(Content,Hash) VALUES(?,?);");err==nil{
            _, err = stmt.Exec(path, source) 
            //w.Write([]byte(fmt.Sprintf("{\"result\":%t}",err==nil)))
        }
    }   
    return nil
}


func up_localmore(w http.ResponseWriter, r *http.Request) {
    //设置内存大小
    r.ParseMultipartForm(32 << 20); //4M
    //获取上传的文件组
    vdir := strings.SplitN(r.URL.Path,"/",3)[2]
    dir := *conf_updir + vdir
    os.MkdirAll(dir, os.ModePerm);//创建上传目录
    files := r.MultipartForm.File["file"]
    var data []string = make([]string,len(files))
    for i,f := range files{
        //打开上传文件
        file, err := f.Open();
        defer file.Close();
        if err != nil {
            log.Info(err);
        }
        //创建上传文件
        cur, err := os.Create(path.Join(dir,f.Filename));
        defer cur.Close();
        if err != nil {
            log.Info(err);
        }
        _, err = io.Copy(cur, file);
        if err != nil {
            fmt.Fprintf(w, "%v", "上传失败")
            return
        }
        data[i] = path.Join(vdir,f.Filename)
        log.Info("上传完成,服务器地址:",data[i])
    }
    responseBody,_ := json.Marshal(map[string]interface{}{"result": 0,"status": "ok","data": data})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseBody)
}




func up_localone(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(32 << 20);
    file, header, err := r.FormFile("file");
    defer file.Close();
    if err != nil {
        log.Fatal(err);
    }
    //创建上传目录
    os.Mkdir("./upload", os.ModePerm);
    //创建上传文件
    cur, err := os.Create("./upload/" + header.Filename);
    defer cur.Close();
    if err != nil {
        log.Fatal(err);
    }
    //把上传文件数据拷贝到我们新建的文件
    io.Copy(cur, file);
}

func up_localmulti(w http.ResponseWriter, r *http.Request) {
    // defer r.Body.Close()  
    // data, _ := ioutil.ReadAll(r.Body) //获取post的数据  
}



func up_ossone(w http.ResponseWriter, r *http.Request) {
    vdir := strings.SplitN(r.URL.Path,"/",3)[2]
    log.Info(vdir)
}