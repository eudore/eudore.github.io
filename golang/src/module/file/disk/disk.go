package disk

import (
	"os"
	"io"
	"fmt"
	"path"
	"net/http"
	"net/url"
	"errors"
	"strings"
	"encoding/json"
	"public/log"
	"module/file/store"
)

type Diskstore struct {
	Config	string
	Host 	string
	Dir 	string
	Key 	string
	Secret	string
}


func (fs *Diskstore) Save(r *http.Request, p string) error {
	r.ParseMultipartForm(32 << 20);
	// var config policy
	// config.debyte = r.MultipartForm.Value["policy"][0]
	// config.check()
	file, _, err := r.FormFile("file");
	// header = 1.png 5668 map[Content-Disposition:[form-data; name="file"; filename="1.png"] Content-Type:[image/png]]
	defer file.Close();
	if err != nil {
		return err
	}
	//创建上传目录
	p = fs.Dir+p
	os.MkdirAll(path.Dir(p), os.ModePerm);
	//创建上传文件
	cur, err := os.Create(p);
	defer cur.Close();
	if err != nil {
		return err
	}
	//把上传文件数据拷贝到我们新建的文件
	io.Copy(cur, file);
	return nil
}


func (fs *Diskstore) Load(w http.ResponseWriter, p string) error {
	file, err := os.Open(fs.Dir+p)
	if err != nil {
		return err
	}
	fileName := url.QueryEscape(path.Base(p)) 
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fileName+"\"")
	_, error := io.Copy(w, file)
	return error
}

func (fs *Diskstore) Callback(r *http.Request) (string, error) {
	return "",nil
}


func (fs *Diskstore) Saves(r *http.Request) ([]string, error) {
	r.ParseMultipartForm(32 << 20);
	//log.Json(r.MultipartForm.Value )
	var config policy
	config.debyte = r.MultipartForm.Value["policy"][0]
	config.check()
	//log.Json(config)


	vdir := strings.SplitN(r.URL.Path,"/",2)[1]
	log.Info("vdir",vdir)
	dir := fs.Dir + vdir
	log.Info(dir)
	log.Info(path.Dir(dir))
	os.MkdirAll(path.Dir(dir), os.ModePerm);//创建上传目录
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
		cur, err := os.Create(dir);
		defer cur.Close();
		if err != nil {
			log.Info(err);
		}
		_, err = io.Copy(cur, file);
		if err != nil {
			//fmt.Fprintf(w, "%v", "上传失败")
			return nil,errors.New(fmt.Sprintf( "%v", "上传失败"))
		}
		data[i] = vdir
		log.Info("上传完成,服务器地址:",vdir)
	}
	return data,nil
}


func New(config string) (store.Store, error) {
	var fs Diskstore
	err := json.Unmarshal([]byte(config), &fs)
	return &fs,err
}


func init() {
	store.Register("disk",New)
}