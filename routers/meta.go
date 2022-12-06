package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"strconv"
)

func Meta(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	params := getparas(req)
	id := params["id"]
	secid := params["secid"]
	mlen := params["len"]
	iid, _ := strconv.Atoi(id)
	isecid, _ := strconv.Atoi(secid)
	ilen, _ := strconv.Atoi(mlen)
	if ilen == 0 {
		ilen = 21
	}
	r, _ := base.Pcontent.GetMeta(iid, isecid, ilen)
	st := sectext{}
	if r != "" {
		st.Text = r
	} else {
		st.Text = "【已结束】"
	}
	json.NewEncoder(w).Encode(st)
}
