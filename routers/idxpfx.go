package routers

import (
	"bytes"
	"fmt"
	"net/http"
	"research/gstr"
	"research/pubgo"
	"research/xbdb"
	"strconv"
	"strings"
)

// 为搜佛说定制。通过关键词查找正反转匹配关键词。
/*
	如：如来藏。为正关键词；
	藏来如。为反转关键词；
*/
func Idxpfx(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	params := getparas(req)
	const (
		tbname   = "c"
		idxfield = "s"
	)

	kw := params["kw"]
	top := "7"
	caid := params["caid"]
	p := params["p"]
	if p == "0" { //表示没有数据了。不过在客户端已经过滤。
		w.Write([]byte(""))
		return
	}
	isr := params["isr"] //是否反转索引.//1是，0或空，否。

	//自动转化参数值
	//idxvalue := Table[string(tbname)].Ifo.FieldChByte(idxfield, kw)
	count := 0
	count, _ = strconv.Atoi(top)
	idx := Newidxfunc(kw, caid, count, isr == "1")
	idx.r.WriteString("[") //idx.r.WriteString("{\"kwpxf\":[")
	if p == "" {           //第一页必须是前缀匹配才能查找。
		if idx.isR {
			rkw := pubgo.Reverse(kw)
			Table["c"].Select.FindPrefixFun([]byte("c"+xbdb.IdxSplit+"r-"+rkw), true, idx.add)
		} else {
			idxvalue := Table[string(tbname)].Ifo.FieldChByte(idxfield, kw)
			Table["c"].Select.WhereIdxLikeFun([]byte(idxfield), idxvalue, true, idx.add)
		}
	} else { //第二页开始，可以根据前一页的最后key定位遍历查找。
		//将pkey的int重新转换回byte，得到正确的k键值
		//这个与search.go的lastkey有点不一样。
		//c.r-一人心念口言-14371||0
		//c.s-一人心念口言-14371||0
		ps := strings.Split(p, xbdb.Split)
		if len(ps) > 1 {
			p = strings.Replace(p, "whfgf", "?", -1) //url转义
			p = strings.Replace(p, "yhfgf", "&", -1) //url转义
			ips := strings.Split(ps[2], idssplit)    //将14371+0转为byte的字符串
			if len(ips) > 1 {
				aid, _ := strconv.Atoi(ips[0])
				sid, _ := strconv.Atoi(ips[1])
				ids := ArtSecToId(aid, sid)
				skey := ps[0] + xbdb.Split + ps[1] + xbdb.Split + ids
				Table[tbname].Select.FindSeekFun([]byte(skey), true, idx.add)
			} else {
				w.Write([]byte("错误的定位页"))
				fmt.Println("错误的定位页p,ps：", p, ps)
				fmt.Println("kw=" + params["kw"])
				return
			}
		}

	}
	jsonstr := idx.r.String()
	if idx.r != nil {
		idx.r.Reset()
		bufpool.Put(idx.r)
	}
	jsonstr = strings.Trim(jsonstr, ",") + "]"
	jsonstr = "{\"Ks\":" + jsonstr + ",\"p\":" + strconv.Quote(idx.pkey) + "}"
	w.Write([]byte(jsonstr))

}

type idxfunc struct {
	kw    string
	caid  string
	count int //返回条数
	loop  int
	isR   bool   //是否反转索引
	pkey  string //当前页最后一个k值，以作下一页的定位////k=c.s-如来-\x00\x00\x19a||\x00\x00\x00?
	r     *bytes.Buffer
}

func Newidxfunc(kw, caid string, count int, isR bool) *idxfunc {
	return &idxfunc{
		kw:    kw,
		caid:  caid,
		count: count,
		isR:   isR,
		r:     bufpool.Get().(*bytes.Buffer),
	}
}
func (i *idxfunc) add(k, v []byte) bool {
	i.pkey = ""
	ks := xbdb.SplitRd(k) //bytes.Split(rd, []byte(xbdb.Split))
	if len(ks) < 3 {
		return true
	}

	//过滤非目录下
	sid := ""
	if i.isR {
		sid = string(ks[3]) //c.r-来如一供士開权--\x00\x00\x1f\"||\x00\x00\x00\x12 (整理数据时出错(--)，多了一个空字段，只能在这里设置一下。)
	} else {
		sid = string(ks[2]) //"c.s-来如一供士開权-\x00\x000\x84||\x00\x00\x00_"
	}
	artid, secid := IdToArtSec(sid) //artid, _ := IdToArtSec(string(ks[1]))
	//artid, secid := IdToArtSec(string(ks[2])) //artid, _ := IdToArtSec(string(ks[1]))
	cid := Artfid[uint32(artid)]
	if !CacaRand(int(cid), i.caid) {
		return true
	}
	var sectext string
	sectext = string(ks[1])
	if i.isR {
		sectext = pubgo.Reverse(string(ks[1])) //反转字符串
	}

	if len(i.kw) == len(sectext) {
		return true
	}
	if !strings.Contains(sectext, i.kw) { //当使用地位遍历的时候，需要进行检验是否已经超过关键词的范围。
		return true
	}
	qsectext := strconv.Quote(sectext)
	if !strings.Contains(qsectext, sectext) { //存在需要转义的，都过滤
		return true
	}
	i.pkey = string(ks[0]) + xbdb.Split + string(ks[1]) + xbdb.Split + strconv.Itoa(artid) + idssplit + strconv.Itoa(secid)
	i.pkey = strings.Replace(i.pkey, "?", "whfgf", -1) //url转义 ?和&不能出现再url。
	i.pkey = strings.Replace(i.pkey, "&", "yhfgf", -1) //url转义 ?和&不能出现再url。

	if !strings.Contains(i.r.String(), qsectext) {
		i.r.WriteString("{\"kw\":" + qsectext + "},")
		//i.r.WriteString("{\"kw\":" + qsectext + ",\"ckey\":" + strconv.Quote(i.pkey) + "},")
		//i.r.WriteString(qsectext + ",")
	} else {
		return true
	}
	i.loop++

	return i.loop < i.count
}
