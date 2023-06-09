package routers

import (
	"bytes"
	"net/http"
	"research/xbdb"
	"strconv"
	"strings"
)

func Meta(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")
	const (
		tbname   = "c"
		idxfield = "id"
	)
	params := getparas(req)
	artid, _ := strconv.Atoi(params["artid"])
	secid, _ := strconv.Atoi(params["secid"])
	id := ArtSecToId(artid, secid)
	//idxvalue := Table[tbname].Select.GetPkKey([]byte(id)) //Table[tbname].Ifo.FieldChByte(idxfield, id)
	bid := xbdb.SplitToCh([]byte(id))
	key := Table[tbname].Select.GetPkKey(bid)
	ef := newmetaexefun(tbname)
	ef.r.Write([]byte("[")) //ef.r.Write([]byte("{\"result\":["))
	Table[tbname].Select.FindSeekFun(key, true, ef.addtext)
	jsonstr := ef.r.String()
	jsonstr = strings.Trim(jsonstr, ",")
	ef.r.Reset()
	ef.r.WriteString(jsonstr)
	ef.r.Write([]byte("]"))
	if ef.r != nil {
		w.Write(ef.r.Bytes())
		//w.Write([]byte(strconv.Quote(ef.r.String()))) //必须使用strconv.Quote转义
		ef.r.Reset()
		bufpool.Put(ef.r)
	}

}

type metaexefun struct {
	r      *bytes.Buffer
	tbname string
	len    int    //最大限制长度
	keys   string // [][]byte
}

func newmetaexefun(tbname string) *metaexefun {
	return &metaexefun{
		r:      bufpool.Get().(*bytes.Buffer),
		tbname: tbname,
		len:    1,
	}
}

// 获取句子并累加，直到大于最大限制长度
func (m *metaexefun) addtext(k, v []byte) bool {
	//rd := xbdb.KVToRd(k, v, []int{})
	m.keys += string(v) //Table[m.tbname].Split(rd) // bytes.Split(rd, []byte(xbdb.Split))
	m.r.WriteString("{\"sec\":" + strconv.Quote(string(v)) + "},")
	//m.r.Write(m.keys[0])
	return len(string(m.keys))*3 < m.len
}
