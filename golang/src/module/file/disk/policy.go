package disk

import (
	"io"
	"fmt"
	"time"
	"hash"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"encoding/base64"
	"public/log"
	"module/file/store"
)

type policy struct {
	Dir		string	`json:"d"`
	Len		int64	`json:"l"`
	Method	string	`json:"m"`
	Expires	int64	`json:"e"`
	debyte	string	`json:"-"`
}


func (p *policy) String() string {
	result,_ := json.Marshal(p)
	p.debyte = base64.StdEncoding.EncodeToString(result)
	return p.debyte
}

func (p *policy) Signed(secret string) string {
	if len(p.debyte) == 0 {
		p.String()
	}
	//debyte := base64.StdEncoding.EncodeToString(p.result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(secret))
	io.WriteString(h, p.debyte)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (p *policy) check() {
	data,_ := base64.StdEncoding.DecodeString(p.debyte)
	json.Unmarshal(data,p)
}




func (s *Diskstore) Signed(p *store.Policy) []byte {
	// create post policy json
	config := policy{
		Dir: p.Directory,
		Len: p.Length,
		Method: p.Method,
		Expires: time.Now().Add(100 * time.Second).Unix(),
	}

	data := make(map[string]interface{})
	data["host"] = s.Host+p.Directory+"?disk"
	data["dir"] = p.Directory
	data["clientID"] = s.Key
	data["policy"] = config.String()
	data["signature"] = config.Signed(s.Secret) 
	//data["callback"] = string(callbackBase64)
	//p.Size = 128<<20
	if p.Length > 4<<20 {
		data["size"] = 1<<20
	}
	response,err :=json.Marshal(data)	
	if err != nil {
		fmt.Println("json err:", err)
	}
	return response
}

func (s *Diskstore) Verify(p *store.Policy) error {
	log.Json(p)
	log.Json(s)
	return nil
}