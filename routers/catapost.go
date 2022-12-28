package routers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func catapost(w http.ResponseWriter, req *http.Request) {

	mu.Lock() //leveldb仅支持单进程数据操作。
	defer mu.Unlock()

	params := postparas(req)
	serpsw := ConfigMap["pws"].(string)                             //服务器端不设置密码，即不可以进行操作
	if serpsw == "" || params["psw"] != ConfigMap["pws"].(string) { //密码不对
		w.Write([]byte("密码不对"))
		return
	}
	//url参数名称必须与表字段名称一致
	//Table["ca"].Insert()
	r := InsOrUpd("ca", params, params["iou"])
	//添加一条以目录标题的空内容文章，以便能够搜索到目录
	if r.Succ {
		paramsrt := map[string]string{
			"id":    "",
			"title": "",
			"fid":   "",
			"split": "。",
			"url":   "",
			"text":  "",
			"psw":   "",
		}
		id := ""
		if params["id"] == "" {
			id = strconv.Itoa(Table["ca"].Ac.GetidNoInc() - 1)
		} else {
			id = params["id"]
		}
		paramsrt["title"] = params["title"]
		paramsrt["fid"] = id
		paramsrt["psw"] = params["psw"]
		InsArt(paramsrt)
	}

	json.NewEncoder(w).Encode(r)
}
