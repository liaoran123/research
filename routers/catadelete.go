package routers

import (
	"encoding/json"
	"net/http"
)

func catadelete(w http.ResponseWriter, req *http.Request) {
	params := postparas(req)
	psw := params["psw"]
	serpsw := ConfigMap["pws"].(string)
	if serpsw == "" || psw != ConfigMap["pws"].(string) { //密码不对
		w.Write([]byte("密码不对"))
		return
	}
	tbname := "ca"
	id := params["id"]
	r := deleteRd(tbname, id)
	json.NewEncoder(w).Encode(r)
}
