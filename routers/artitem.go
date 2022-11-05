package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"research/gstr"
	"research/pubgo"
	"strconv"
)

func Artitem(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	params := getparas(req)

	id := params["id"]
	iid, _ := strconv.Atoi(id)
	count := params["count"]
	icount, _ := strconv.Atoi(count)
	p := params["p"] //当前页

	//art := base.NewArticle()
	r := base.PArticle.GetCataArt(iid, icount, p)
	json.NewEncoder(w).Encode(r)

}
