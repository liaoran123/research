package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"strconv"
)

//获取一文章信息
func artget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req)
	id := params["id"]
	secid := params["secid"]
	if secid != "" { //参数有secid，即是//获取文章一段落内容artgetsec
		artgetsec(w, req)
		return
	}
	iid, _ := strconv.Atoi(id)
	r := base.Pcontent.GetArtInfo(iid)
	json.NewEncoder(w).Encode(r)

}
