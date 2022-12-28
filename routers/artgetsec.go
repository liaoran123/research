package routers

import (
	"net/http"
	"research/xbdb"
	"strconv"
)

//获取文章一段落内容
func artgetsec(w http.ResponseWriter, req *http.Request) {
	const (
		tbname   = "c"
		idxfield = "id"
	)
	params := getparas(req)
	id := params["id"]
	secid := params["secid"]
	//文章id+句子段落id
	iid, _ := strconv.Atoi(id)
	bid := xbdb.IntToBytes(iid) //必须这样转换，否则排序不正确
	aid := string(bid) + "+"
	iid, _ = strconv.Atoi(secid)
	bid = xbdb.IntToBytes(iid)
	aid += string(bid) //+ "+"
	getonerecord(tbname, idxfield, id, w)
}

type sectext struct {
	Text string `json:"text" `
}
