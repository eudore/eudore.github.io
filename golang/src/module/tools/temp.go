package tools

import (
	"io"
	"html/template"
)

var templates map[string]*template.Template

func Template(wr io.Writer,sour string, data interface{}) (err error) {
	tmp, ok := templates[sour]
	if !ok {
		tmp,err = template.ParseFiles("/data/web/templates/"+sour,"/data/web/templates/base.html")
		if err != nil{
			return
		}
		// no cache
		//templates[sour]=tmp
	}
	return tmp.Execute(wr, data)
}

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
}