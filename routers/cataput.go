package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"research/base"
	"research/pubgo"
)

func cataput(w http.ResponseWriter, req *http.Request) {

	//ts := pubgo.Newts() //计算执行时间

	params := postparas(req)
	Rmsg := NewRmsg()

	id := params["id"]
	iid, _ := strconv.Atoi(id)
	fid := params["fid"]
	ifid, _ := strconv.Atoi(fid)
	name := params["name"]
	psw := params["psw"]
	isleaf := params["isleaf"]
	if psw != pubgo.ConfigMap["pws"].(string) { //密码不对
		Rmsg.Msg = "密码不对"
		json.NewEncoder(w).Encode(Rmsg)
		return
	}
	r := base.Pcata.Put(iid, ifid, name, isleaf)
	//ys := ts.Gstrts()
	Rmsg.Succ = r
	if r {
		Rmsg.Msg = "提交成功。"
	} else {
		Rmsg.Msg = "提交失败。"
	}
	//Rmsg.Time = ys
	json.NewEncoder(w).Encode(Rmsg)
}
