package routers

import (
	"html/template"
	"net/http"
)

func Searchhtml(w http.ResponseWriter, req *http.Request) {
	var t *template.Template
	t, _ = template.ParseFiles("admin/search.html") //从文件创建一个模板
	t.Execute(w, nil)
}
