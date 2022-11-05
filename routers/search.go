package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"research/pubgo"
)

func Search(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")
	//ts := pubgo.Newts() //计算执行时间
	pubgo.Tj.Brows("Search")
	params := getparas(req)
	Rmsg := NewRmsg()

	kw := params["kw"]
	p := params["p"] //当前页
	count := params["count"]
	caids := params["caids"] //目录id,在目录范围下搜索.支持多个，用 “|”隔开。
	if kw == "" {
		Rmsg.Msg = "请输入搜索词！"
		json.NewEncoder(w).Encode(Rmsg)
		return
	}
	order := params["order"]
	r := base.PSe.Search(kw, p, count, caids, order == "0")
	//ys := ts.Gstrts()
	json.NewEncoder(w).Encode(r)
	r.Reset() //置空
	base.RsetAllPool.Put(r)
}
