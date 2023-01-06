package routers

import (
	"net/http"
	"research/gstr"
	"research/pubgo"
	"research/xbdb"
	"strconv"
)

func Artitem(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))
	const (
		tbname   = "art"
		idxfield = "fid"
	)
	params := getparas(req)

	idxvalue := Table[string(tbname)].Ifo.FieldChByte(idxfield, params["id"])
	count := -1
	if params["count"] != "" {
		count, _ = strconv.Atoi(params["count"])
	}
	b := 0
	if params["b"] != "" {
		b, _ = strconv.Atoi(params["b"])
	}
	tbd := Table[tbname].Select.WhereIdx([]byte(idxfield), idxvalue, true, b, count)
	if tbd == nil {
		return
	}
	r := Table[tbname].DataToJson(tbd) // DataToJson(tbd, Table[tbname].Ifo)
	if r != nil {
		w.Write(r.Bytes())
		//w.Write([]byte(strconv.Quote(r.String()))) //必须使用strconv.Quote转义
		r.Reset()
		xbdb.Bufpool.Put(r)
	}

}
