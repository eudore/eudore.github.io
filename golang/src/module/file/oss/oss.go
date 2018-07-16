//  by  https://www.alibabacloud.com/help/zh/doc-detail/31927.htm?spm=a3c0i.o31926zh.b99.82.435c6eb4bRqYJV#%E4%BB%A3%E7%A0%81%E4%B8%8B%E8%BD%BD

/*
    Browser:    https://help.aliyun.com/document_detail/31927.html?spm=a2c4g.11186623.6.634.91y8Mi#%E4%BB%A3%E7%A0%81%E4%B8%8B%E8%BD%BD
    Policy:     https://help.aliyun.com/document_detail/31927.htm?spm=a3c0i.o31926zh.b99.82.435c6eb4bRqYJV#%E4%BB%A3%E7%A0%81%E4%B8%8B%E8%BD%BD
    Callback:   https://help.aliyun.com/document_detail/50092.html?spm=a2c4g.11186623.6.1089.kGyEEu#%E8%B0%83%E8%AF%95%E5%9B%9E%E8%B0%83%E6%9C%8D%E5%8A%A1%E5%99%A8
*/

package oss;

import (
	"net/http"
	"encoding/json"
	"module/file/store"
)

type Ossstore struct {
	Config	string
	bucket	string
	Host 	string
	Dir		string
	Key 	string
	Secret	string
}


func (s *Ossstore) Load(w http.ResponseWriter, p string) error {
	http.Redirect(w,r,s.Host + s.Dir + p,302)
	return nil
}

func New(config string) (store.Store, error) {
	var oss Ossstore
	err := json.Unmarshal([]byte(config), &oss)
	oss.Config=config
	return &oss,err
}

func init() {
	store.Register("oss",New)
}