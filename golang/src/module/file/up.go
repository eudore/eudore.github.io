package file;

import (
	"strconv"
	"net/http"
	"encoding/json"
	"public/router"
	"public/log"
	"module/file/store"
)
func fileup(w http.ResponseWriter, r *http.Request) {
	user := router.GetValue(r, "user")
	zone := router.GetValue(r, "zone")
	log.Info(r.URL.Path)
	fs,err := store.Getstore(user+"/"+zone)
	log.Json(fs)
	if err!=nil {
		w.WriteHeader(http.StatusNotFound) 
		return
	}
	// get policy info
/*	if r.TLS == nil {
		p.Host = "http://"+r.Host
	}*/
	// signed
/*	if r.URL.RawQuery == "" {
		response := fs.Signed(p)
		log.Json(p)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}*/
		// if(true || strings.HasPrefix(r.Header.Get("Referer"),"https://www.wejass.com/file/")){

	name := r.URL.Path
	if r.URL.RawQuery == "callback" {
		// callback
		log.Info("callback file")
		name,err = fs.Callback(r)
	}else {
/*		// verify
		p := &store.Policy{
			Host: "https://"+r.Host,
			Directory: "/file"+r.URL.Path,
			Method: "POST",
			Length: 0,
		}
		log.Info("Verify")
		if fs.Verify(p) != nil {
			w.WriteHeader(http.StatusOK)
			return
		}*/
		// upload
		log.Info("save file")
		err = fs.Save(r,name)
	}
	// add file and return info
	if err == nil{
		store.Add(user+"/"+zone,[]string{name})
		responseBody,_ := json.Marshal(map[string]interface{}{"status":"ok","data": Prefix + name})
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
		log.Info("Post Response : 200 OK . uri: ",name)   
	}else {
		log.Info(err)
	}
}
