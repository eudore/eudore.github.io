package tools

import (
    "crypto/md5"
    "encoding/base64"
	"io/ioutil"
    "encoding/hex"
)

func Md5(str string) (string){
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(str))
    return hex.EncodeToString(md5Ctx.Sum(nil))
}


func Base64(data string) (string){
    return base64.StdEncoding.EncodeToString([]byte(data))
}

func Md5SumFile(file string) (value [md5.Size]byte, err error) {
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return
    }
    value = md5.Sum(data)
    return
}
