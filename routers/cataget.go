package routers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"research/base"
)

var catamap map[string]func(string, http.ResponseWriter)

func cataget(w http.ResponseWriter, req *http.Request) {
	if catamap == nil {
		catamap = make(map[string]func(string, http.ResponseWriter))
		catamap["id"] = getcatainfo
		catamap["fid"] = getChildCatas
		catamap["dir"] = getCataDir
	}

	params := getparas(req)
	for k, v := range params {
		if f, ok := catamap[k]; ok {
			f(v, w)
			return
		}
	}
}

//返回一目录信息
func getcatainfo(id string, w http.ResponseWriter) {
	iid, _ := strconv.Atoi(id)
	r := base.Pcata.GetCata(iid)
	json.NewEncoder(w).Encode(r)
}

//返回所有子目录信息
func getChildCatas(id string, w http.ResponseWriter) {
	ifid, _ := strconv.Atoi(id)
	r := base.Pcata.ChildCatas(ifid)
	json.NewEncoder(w).Encode(r)
}

//获取目录路径
func getCataDir(id string, w http.ResponseWriter) {
	idir, _ := strconv.Atoi(id)
	r := base.CRAMs.GetCataDir(idir)
	json.NewEncoder(w).Encode(r)
}
