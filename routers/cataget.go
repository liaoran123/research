package routers

import (
	"net/http"
	"strconv"
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
	const (
		tbname   = "ca"
		idxfield = "id"
	)
	getonerecord(tbname, idxfield, id, w)
}

//返回所有子目录信息
func getChildCatas(id string, w http.ResponseWriter) {
	const (
		tbname   = "ca"
		idxfield = "fid"
	)
	idxvalue := Table[string(tbname)].Ifo.FieldChByte(idxfield, id)
	tbd := Table[string(tbname)].Select.WhereIdx([]byte(idxfield), idxvalue, 0, -1)
	if tbd == nil {
		return
	}
	r := DataToJson(tbd, Table[string(tbname)].Ifo)
	if r != nil {
		w.Write(r.Bytes())
		//w.Write([]byte(strconv.Quote(r.String()))) //必须使用strconv.Quote转义
		r.Reset()
		bufpool.Put(r)
	}

}

//获取目录路径
func getCataDir(id string, w http.ResponseWriter) {
	idir, _ := strconv.Atoi(id)
	r := CRAMs.GetCaDirToJson(idir) // .GetCataDir(idir)
	w.Write([]byte(r))
	//w.Write([]byte(strconv.Quote(Se.r.String()))) //必须使用strconv.Quote转义
	//json.NewEncoder(w).Encode(r)  strconv.Unquote
}
