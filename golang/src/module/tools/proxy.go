package tools;

import (
	"fmt"
	"strings"
	"bytes"
    "net/http"
	"html/template"
    "encoding/base64"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Proxy struct {
	Server string
	Port int
	Protocol string
	Method string
	Obfs string
	Password string
	Obfsparam string
	Protoparam string
	Remarks string
	Group string
}

func (p Proxy) String() string {
	var data string = fmt.Sprintf("%s:%d:%s:%s:%s:%s/?obfsparam=%s&protoparam=%s&remarks=%s&group=%s", p.Server, p.Port,p.Protocol,p.Method,p.Obfs, encode(p.Password),encode(p.Protoparam),encode(p.Obfsparam),encode(p.Remarks),encode(p.Group) )
    return fmt.Sprintf("ssr://%s",encode(data))
}

func (p Proxy) Html() string {
	return fmt.Sprintf("<tr id=\"\"><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", p.Server, p.Port,p.Protocol,p.Method,p.Obfs, p.Password,p.Protoparam,p.Obfsparam,p.Remarks,p.Group)
}




func proxy(w http.ResponseWriter, r *http.Request) {	
	var doc bytes.Buffer
	var data = []Proxy{};
	tmp, err := template.ParseFiles("/data/web/templates/tools/proxy.html","/data/web/templates/base.html")
    if err == nil {
        tmp.Execute(&doc,map[string]interface{}{"data": data})
		w.Write([]byte(doc.String()))
    }
}
func subscribe(w http.ResponseWriter, r *http.Request) {	
	host := r.Header.Get("X-Real-Ip");
	if host=="" {
		host = strings.Split(r.RemoteAddr,":")[0]
	}
	if host != "176.122.165.113" {
		w.WriteHeader(404)
        return
	}

	var data []string;
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		rows, _ := db.Query("SELECT `Server`,`Port`,`Protocol`,`Method`,`Obfs`,`Password`,`Obfsparam`,`Protoparam`,`Remarks`,`Group` FROM tb_tools_proxy WHERE `Enable`=true ORDER BY `Remarks`;")
		for rows.Next() {
    		line := Proxy{}
    		rows.Scan(&line.Server,&line.Port,&line.Protocol,&line.Method,&line.Obfs,&line.Password,&line.Obfsparam,&line.Protoparam,&line.Remarks,&line.Group)
    		data = append(data, line.String())
		}
	}
	w.Write([]byte(encode(strings.Join(data,"\n"))))
}

func encode(data string) (string){
    return base64.RawURLEncoding.EncodeToString([]byte(data))
}