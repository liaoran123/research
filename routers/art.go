package routers

import (
	"net/http"
	"research/pubgo"
)

//var mu sync.RWMutex

var (
	artmethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。
)

func Art(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/art/" + req.Method)

	//req.artmethod
	if artmethod == nil {
		artmethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		artmethod["POST"] = artpost     //添加
		artmethod["GET"] = artget       //打开
		artmethod["DELETE"] = artdelete //删除
		artmethod["PUT"] = artput       //修改
	}
	if f, ok := artmethod[req.Method]; ok {
		f(w, req)
	}
}
