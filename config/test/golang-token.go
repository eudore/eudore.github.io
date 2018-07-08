package main


import (
	"fmt"
	"strings"
	"encoding/json"
	"encoding/base64"
	"encoding/hex"
	"crypto/sha256"
	"crypto/hmac"
	"golang.org/x/crypto/bcrypt"
)

// SALT 密钥
const SALT = "secret"

// Header 消息头部
type Header struct {
	Typ string // Token Type
	Alg string // Message Authentication Code Algorithm - The issuer can freely set an algorithm to verify the signature on the token. However, some asymmetrical algorithms pose security concerns
//	Cty string // Content Type - This claim should always be JWT
}

// PayLoad 负载
type PayLoad struct {
	Sub string `json:"sub"`
	Name string `json:"name"`
	Admin bool `json:"admin"`
	// Expire int `json:"exp"`
}


// JWT 完整的本体
type JWT struct {
	Header				`json:"header"`
	PayLoad				`json:"payload"`
	Signature	string	`json:"signature"`
}

var tk string
func main(){
	// Encode
	jwt := JWT{}
	jwt.Header = Header{"JWT","HS256"}
	jwt.PayLoad = PayLoad{"1234567890","John Doe",true}
	result := jwt.Encode()
	tk = result
	fmt.Println(result)

	//Decode
	testStr := tk
	jwt = JWT{}
	if jwt.Decode(testStr) {
		fmt.Println(jwt)
	} else {
		fmt.Println("error json content")
	}
}

func getHmacCode(s string) string {
    h := hmac.New(sha256.New, []byte(SALT))
	h.Write([]byte(s))
	key := h.Sum(nil)
    return hex.EncodeToString(key)
}

func getBcryptCode(s string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt. DefaultCost)
    return hex.EncodeToString(hash)
}


func (jwt *JWT) Encode() string {
	header, err := json.Marshal(jwt.Header)
	checkError(err)
	headerString := base64.StdEncoding.EncodeToString(header)
	payload, err := json.Marshal(jwt.PayLoad)
	payloadString := base64.StdEncoding.EncodeToString(payload)
	checkError(err)
	
	format := headerString + "." + payloadString
    signature := getHmacCode(format)

	return format + "." + signature
}


// Decode 验证 jwt 签名是否正确,并将json内容解析出来
func (jwt *JWT) Decode( code string) bool {

	arr := strings.Split(code,".")
	if len(arr) != 3 {
		return false
	}

	// 验证签名是否正确
	format := arr[0] + "." + arr[1]
	signature := getHmacCode(format)
	if signature != arr[2] {
		return false
	}


	header, err := base64.StdEncoding.DecodeString(arr[0])
	checkError(err)
	payload, err := base64.StdEncoding.DecodeString(arr[1])
	checkError(err)

	json.Unmarshal(header, &jwt.Header)
	json.Unmarshal(payload,&jwt.PayLoad)

	return true
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
