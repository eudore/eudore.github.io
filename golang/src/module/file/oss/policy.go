package oss;

import (
	"io"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
	"hash"
	"net/url"
	"encoding/json"
	"encoding/base64"
	"module/file/store"
)


const (
	expire_time =   60
	callbackBody=   `{"filename":${object},"mimeType":${mimeType},"size":${size}}`
	callbackBodyType="application/json"
	base64Table =   "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func getGmtIso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire 
}

type ConfigStruct struct{
	Expiration string `json:"expiration"`
	Conditions [][]interface{} `json:"conditions"`
} 

type CallbackParam struct{
	CallbackUrl string `json:"callbackUrl"`
	CallbackBody string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}


func (s *Ossstore) Signed(p *store.Policy) []byte {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = getGmtIso8601(expire_end)

	// create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire  
	if p.Directory!=""{
		p.Directory = p.Directory[1:]
		var conditiondir []interface{}
		if p.Directory[len(p.Directory):]=="/" {
			conditiondir = append(conditiondir, "starts-with")
			conditiondir = append(conditiondir, "$key")
			conditiondir = append(conditiondir, p.Directory)
		}else{
			conditiondir = append(conditiondir, "eq")
			conditiondir = append(conditiondir, "$key")
			conditiondir = append(conditiondir, p.Directory)
				
		}
		config.Conditions = append(config.Conditions, conditiondir)
	}
	if p.Length!=0 {
		var conditionlen []interface{}
		conditionlen = append(conditionlen, "content-length-range")
		conditionlen = append(conditionlen, p.Length)
		conditionlen = append(conditionlen, p.Length)
		config.Conditions = append(config.Conditions, conditionlen)
	}	
	indent,_ := json.MarshalIndent(&config, "", "\t")
	fmt.Println(string(indent))

	// calucate signature
	result,err:=json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(s.Secret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// callback args
	var callbackParam CallbackParam
	callbackParam.CallbackUrl = p.Host + "/" + url.QueryEscape(p.Directory)+"?callback"
	callbackParam.CallbackBody = callbackBody
	callbackParam.CallbackBodyType = callbackBodyType	
	indent,_ = json.MarshalIndent(&callbackParam, "", "\t")
	fmt.Println(string(indent))
	callback_str,err:=json.Marshal(callbackParam)
	if err != nil {
		fmt.Println("callback json err:", err)
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	data := make(map[string]interface{})
	data["host"] = s.Host
	data["key"] = p.Directory
	data["OSSAccessKeyId"] = s.Key
	data["signature"] = string(signedStr)
	data["policy"] = string(debyte)
	data["callback"] = string(callbackBase64)
	data["success_action_status"]=200
	response,err:=json.Marshal(data)
	if err != nil {
		fmt.Println("json err:", err)
	}
	return response

}

func (s *Ossstore) Verify(*store.Policy) error {
	return nil
}
