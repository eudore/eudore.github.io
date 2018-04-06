//  by  https://www.alibabacloud.com/help/zh/doc-detail/31927.htm?spm=a3c0i.o31926zh.b99.82.435c6eb4bRqYJV#%E4%BB%A3%E7%A0%81%E4%B8%8B%E8%BD%BD


/*
    Browser:    https://help.aliyun.com/document_detail/31927.html?spm=a2c4g.11186623.6.634.91y8Mi#%E4%BB%A3%E7%A0%81%E4%B8%8B%E8%BD%BD
    Policy:     https://help.aliyun.com/document_detail/31927.htm?spm=a3c0i.o31926zh.b99.82.435c6eb4bRqYJV#%E4%BB%A3%E7%A0%81%E4%B8%8B%E8%BD%BD
    Callback:   https://help.aliyun.com/document_detail/50092.html?spm=a2c4g.11186623.6.1089.kGyEEu#%E8%B0%83%E8%AF%95%E5%9B%9E%E8%B0%83%E6%9C%8D%E5%8A%A1%E5%99%A8
*/

package file;

import (
	"io"
	"net/http"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
	"encoding/json"
	"hash"
	"encoding/base64"
)



var coder = base64.NewEncoding(base64Table)
func base64Encode(src []byte) []byte {
    return []byte(coder.EncodeToString(src))
}

func get_gmt_iso8601(expire_end int64) string {
    var tokenExpire = time.Unix(expire_end, 0).Format("2006-01-02T15:04:05Z")
    return tokenExpire 
}

type ConfigStruct struct{
    Expiration string `json:"expiration"`
    Conditions [][]string `json:"conditions"`
} 

type PolicyToken struct{
    AccessKeyId string `json:"accessid"`
    Host string `json:"host"`
    Expire int64 `json:"expire"`
    Signature string `json:"signature"`
    Policy string `json:"policy"`
    Directory string `json:"dir"`
    Callback string `json:"callback"`
}

type CallbackParam struct{
    CallbackUrl string `json:"callbackUrl"`
    CallbackBody string `json:"callbackBody"`
    CallbackBodyType string `json:"callbackBodyType"`
}

func get_policy_token() []byte {
    now := time.Now().Unix()
    expire_end := now + expire_time
    var tokenExpire = get_gmt_iso8601(expire_end)

    //create post policy json
    var config ConfigStruct
    config.Expiration = tokenExpire  
    var condition []string
    condition = append(condition, "starts-with")
    condition = append(condition, "$key")
    condition = append(condition, *conf_upload_dir)
    config.Conditions = append(config.Conditions, condition)

    //calucate signature
    result,err:=json.Marshal(config)
    debyte := base64.StdEncoding.EncodeToString(result)
    h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(*conf_accessKeySecret))
    io.WriteString(h, debyte)
    signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

    var callbackParam CallbackParam
    callbackParam.CallbackUrl = callbackUrl
    callbackParam.CallbackBody = callbackBody//"filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}"
    callbackParam.CallbackBodyType = callbackBodyType//"application/x-www-form-urlencoded"
    callback_str,err:=json.Marshal(callbackParam)
    if err != nil {
        fmt.Println("callback json err:", err)
    }
    callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

    var policyToken PolicyToken
    policyToken.AccessKeyId = *conf_accessKeyId
    policyToken.Host = *conf_host
    policyToken.Expire = expire_end
    policyToken.Signature = string(signedStr)
    policyToken.Directory = *conf_upload_dir
    policyToken.Policy = string(debyte)
    policyToken.Callback = string(callbackBase64)
    response,err:=json.Marshal(policyToken)
    if err != nil {
        fmt.Println("json err:", err)
    }
    return response
}

func oss_policy(w http.ResponseWriter, r *http.Request) {
	response := get_policy_token()
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(response)
}