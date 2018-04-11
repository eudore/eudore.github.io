package file;

import (
    "fmt"
    "os"
    "io"
    "path"
    "strings"
    "net/http"
    "encoding/json"
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
    uptype["localpart"]     =   up_localpart
    uptype["ossone"]        =   up_localmulti
    uptype["ossmulti"]      =   up_localmulti
    uptype["osscallback"]   =   up_localmulti
}

func file_up(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    if v,ok := uptype[r.Form["type"][0]]; ok {
        v(w,r)
    }else{
        log.Warning("invalid upload file type",r.Form["type"])
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
