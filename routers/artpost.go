package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"research/pubgo"
	"strconv"
)

func artpost(w http.ResponseWriter, req *http.Request) {
	//ts := pubgo.Newts() //计算执行时间

	params := postparas(req)
	Rmsg := NewRmsg()

	id := params["id"]
	title := params["title"]
	text := params["text"]
	//text = strings.Replace(text, "-", "﹣", -1) //-是系统保留字，需要转义为﹣。
	/*
		isdir := params["isdir"]
			if isdir == "" {                           ////默认是文章
				isdir = "0"
			}
			if isdir != "0" { //以防填错纠正。
				isdir = "1"
			}*/
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

	ifid, _ := strconv.Atoi(fid)
	iid, _ := strconv.Atoi(id)
	r := base.PArticle.Insert(iid, title, text, split, url, ifid)
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
