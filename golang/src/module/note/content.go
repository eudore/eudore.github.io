package note

import (
	"net/http"
	"public/log"
	"module/global"
	//"module/tools"
)


func getnote(w http.ResponseWriter, r *http.Request) {
	n := NewNote(r.URL.Path)
	n.LoadData()
	if len(n.EditTime)!= 0 {
		n.Show()
		err := global.Template(w,"note/content.html",map[string]interface{}{"Uri": r.RequestURI,"Note": n,"Level": 0})
		if err != nil {
			log.Info(err)	
		}
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}


func upnote(w http.ResponseWriter, r *http.Request) {
}