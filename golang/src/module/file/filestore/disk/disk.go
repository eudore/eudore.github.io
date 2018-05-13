package disk

import (
	"os"
	"io"
	"fmt"
	"path"
	"net/http"
	"errors"
	"strings"
	"encoding/json"
	"public/log"
	"module/file/filestore"
)

type Diskstore struct {
	Host 	string
	Dir 	string
}

func (fs *Diskstore) Policy(p *filestore.PolicyInfo) []byte {
	data := make(map[string]interface{})
	data["host"] = fs.Host+p.Directory+"?disk"
	if p.Length > 4<<20 {
		data["size"] = 1<<20
	}
	//p.Size = 128<<20
	response,_ :=json.Marshal(data)
	return response
}

func (fs *Diskstore) Save(r *http.Request) ([]string, error) {
	r.ParseMultipartForm(32 << 20);
	vdir := strings.SplitN(r.URL.Path,"/",3)[2]
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



func New(config string) (filestore.Store, error) {
	var fs Diskstore
	err := json.Unmarshal([]byte(config), &fs)
	return &fs,err
}


func init() {
	filestore.Register("disk",New)
}
/*
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
*/









/*
func downloadFile(fileFullPath string, res *restful.Response) {
	file, err := os.Open(fileFullPath)

	if err != nil {
		res.WriteEntity(_dto.ErrorDto{Err: err})
		return
	}

	defer file.Close()
	fileName := path.Base(fileFullPath)
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	res.AddHeader("Content-Type", "application/octet-stream")
	res.AddHeader("content-disposition", "attachment; filename=\""+fileName+"\"")
	_, error := io.Copy(res.ResponseWriter, file)
	if error != nil {
		res.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}*/