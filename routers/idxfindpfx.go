package routers

import (
	"bytes"
	"net/http"
	"research/gstr"
	"research/pubgo"
	"research/xbdb"
	"strconv"
	"strings"
)

func Idxfindpfx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	params := getparas(req)
	const (
		tbname   = "c"
		idxfield = "s"
	)

	kw := params["kw"]
	top := params["top"]
	if top == "" {
		top = "21"
	}
	caid := params["caid"]

	//自动转化参数值
	idxvalue := Table[string(tbname)].Ifo.FieldChByte(idxfield, kw)
	count := 0
	count, _ = strconv.Atoi(top)
	idx := NewidxExefunc(kw, caid, count)
	idx.r.WriteString("[")
	Table["c"].Select.WhereIdxLikeFun([]byte(idxfield), idxvalue, true, idx.add)
	jsonstr := idx.r.String()
	if idx.r != nil {
		idx.r.Reset()
		bufpool.Put(idx.r)
	}
	jsonstr = strings.Trim(jsonstr, ",") + "]"
	w.Write([]byte(jsonstr))
	//w.Write([]byte(strconv.Quote(ef.r.String()))) //必须使用strconv.Quote转义

}

type idxExefunc struct {
	kw    string
	caid  string
	count int //返回条数
	loop  int
	r     *bytes.Buffer
}

func NewidxExefunc(kw, caid string, count int) *idxExefunc {
	return &idxExefunc{
		kw:    kw,
		caid:  caid,
		count: count,
		r:     bufpool.Get().(*bytes.Buffer),
	}
}
func (i *idxExefunc) add(rd []byte) bool {
	ks := bytes.Split(rd, []byte(xbdb.Split))
	//过滤非目录下
	artid, _ := IdToArtSec(string(ks[1]))
	cid := Artfid[uint32(artid)]
	if !CacaRand(int(cid), i.caid) {
		return true
	}
	if strings.Contains(string(ks[0]), " ") {
		return true
	}
	if strings.Contains(string(ks[0]), "\u3000") {
		return true
	}
	if !strings.Contains(i.r.String(), string(ks[0])) {
		i.r.WriteString("{\"key\":\"" + strconv.Quote(string(ks[0])) + "\"},")
	} else {
		return true
	}
	i.loop++
	return i.loop < i.count
}
