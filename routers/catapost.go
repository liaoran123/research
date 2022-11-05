package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"research/pubgo"
	"strconv"
)

func catapost(w http.ResponseWriter, req *http.Request) {

	//ts := pubgo.Newts() //计算执行时间

	//mu.Lock() //leveldb仅支持单进程数据操作。
	//defer mu.Unlock()

	params := postparas(req)
	Rmsg := NewRmsg()
	id := params["id"]
	title := params["title"]
	fid := params["fid"]
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

	r := base.Pcata.Insert(iid, title, ifid)
	//ys := ts.Gstrts()
	if r {
		Rmsg.Msg = "提交成功"
		Rmsg.Succ = true
		//Rmsg.Time = ys
	}else{
		Rmsg.Msg = "提交失败。"
	}
	json.NewEncoder(w).Encode(Rmsg)
}
