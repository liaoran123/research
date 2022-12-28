package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"research/pubgo"
	"strconv"
)

func artput(w http.ResponseWriter, req *http.Request) {
	//ts := pubgo.Newts() //计算执行时间

	params := postparas(req)
	Rmsg := NewRmsg()

	id := params["id"]
	title := params["title"]

	text := params["text"]

	fid := params["fid"]
	if fid == "" {
		fid = "0" //默认是0，即顶级目录
	}
	split := params["split"]
	url := params["url"]
	psw := params["psw"]
	if psw != pubgo.ConfigMap["pws"].(string) { //密码不对
		Rmsg.Msg = "密码不对"
		json.NewEncoder(w).Encode(Rmsg)
		return
	}
	if title == "" {
		Rmsg.Msg = "标题不能为空"
		json.NewEncoder(w).Encode(Rmsg)
		return
	}
	iid, _ := strconv.Atoi(id)
	ifid, _ := strconv.Atoi(fid)

	r := base.PArticle.Put(iid, title, text, split, url, ifid)
	//ys := ts.Gstrts()
	if r {
		Rmsg.Msg = "提交成功"
		Rmsg.Succ = true
		//Rmsg.Time = ys
	} else {
		Rmsg.Msg = "提交失败。"
	}
	json.NewEncoder(w).Encode(Rmsg)
}
