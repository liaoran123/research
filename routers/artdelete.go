package routers

import (
	"encoding/json"
	"net/http"
	"research/xbdb"
	"strconv"
)

func artdelete(w http.ResponseWriter, req *http.Request) {
	params := postparas(req)
	Rmsg := NewRmsg()
	serpsw := ConfigMap["pws"].(string)                             //服务器端不设置密码，即不可以进行操作
	if serpsw == "" || params["psw"] != ConfigMap["pws"].(string) { //密码不对
		Rmsg.Msg = "密码不对"
		json.NewEncoder(w).Encode(Rmsg)
		return
	}
	id := params["id"]
	r := deleteRd("art", id)
	if !r.Succ {
		json.NewEncoder(w).Encode(r)
		return
	}
	iid, _ := strconv.Atoi(params["id"])
	bid := xbdb.IntToBytes(iid)
	idxvalue := Table["c"].Ifo.FieldChByte("id", string(bid))
	df := new(delfun)
	Table["c"].Select.WherePKLikeFun(idxvalue, 0, -1, true, df.delc)
	json.NewEncoder(w).Encode(df.r)
}

type delfun struct {
	r xbdb.ReInfo
}

func (d *delfun) delc(k, v []byte) bool {
	//ks := bytes.Split(rd, []byte(xbdb.Split))
	//d.r = Table["c"].Delete(ks[0])
	d.r = Table["c"].Delete(k) //.Deletekey(k)
	return d.r.Succ
}
