package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"research/gstr"
	"research/pubgo"
	"strconv"
)

func Idxfindpfx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	params := getparas(req)

	kw := params["kw"]
	top := params["top"]
	itop, _ := strconv.Atoi(top)
	r := base.Pcontent.Idx.GetPfx(kw, itop)
	ip := Idxpfx{}
	ip.Pfx = r
	json.NewEncoder(w).Encode(ip)
}

type Idxpfx struct {
	Pfx []string `json:"pfx"`
}
