package routers

import (
	"net/http"
	"research/gstr"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func NewEngine() *Engine { //静态变量返回值用指针*,非静态类则返回copy一份
	return &Engine{router: make(map[string]HandlerFunc)}
}
func (engine *Engine) Addrouter(path string, hf HandlerFunc) {
	engine.router[path] = hf
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//-以第一个路径目录为路由节点
	pathkey := gstr.Mstr(req.URL.Path, "/", "/")
	if pathkey == "" {
		pathkey = "/"
	} else {
		pathkey = "/" + pathkey + "/"
	}
	handler, ok := engine.router[pathkey]
	if ok {
		handler(w, req)
	} else {
		handler, ok = engine.router["/redir/"] //跳转错误页
		if ok {
			handler(w, req)
		}
	}
}
