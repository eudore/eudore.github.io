package file;

import (
    "fmt"
    "os"
    "io"
    "strings"
    "net/http"
    "path"
    "encoding/json"
    "public/log"
)

func file_up(w http.ResponseWriter, r *http.Request) {
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
    responseBody,_ := json.Marshal(map[string]interface{}{"status":"ok","data":data})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseBody)
}


func file_split(w http.ResponseWriter, r *http.Request) {
    // defer r.Body.Close()  
    // data, _ := ioutil.ReadAll(r.Body) //获取post的数据  
}


func uploadOne(w http.ResponseWriter, r *http.Request) {
    //设置内存大小
    r.ParseMultipartForm(32 << 20);
    //获取上传的第一个文件
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