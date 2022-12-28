package routers

import (
	"net/http"
	"research/pubgo"
)

var (
	catamethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

func Cata(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/Cata/" + req.Method)

	//req.catamethod
	if catamethod == nil {
		catamethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		catamethod["POST"] = catapost     //添加
		catamethod["GET"] = cataget       //查询
		catamethod["DELETE"] = catadelete //删除
		catamethod["PUT"] = catapost      //cataput       //修改
	}
	if f, ok := catamethod[req.Method]; ok {
		f(w, req)
	}
}
