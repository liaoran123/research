package routers

import (
	"bytes"
	"net/http"
	"research/xbdb"
	"strconv"
)

//获取一文章信息
func artget(w http.ResponseWriter, req *http.Request) {
	const (
		tbname   = "art"
		idxfield = "id"
	)
	params := getparas(req)
	//getonerecord(tbname, idxfield, params["id"], w)
	//打开文章表
	idxvalue := Table[tbname].Ifo.FieldChByte(idxfield, params["id"])
	tdb := Table[tbname].Select.OneRecord(idxvalue)
	if tdb == nil {
		return
	}
	//打开文章内容表，一个文章表对应多个内容表
	cf := newcexefun()
	//以下的id转换比较复杂，对应添加时的转换ArtSecToId函数
	id, _ := strconv.Atoi(params["id"])
	bid := xbdb.IntToBytes(id)
	idxvalue = Table["c"].Ifo.FieldChByte(idxfield, string(bid))
	Table["c"].Select.WherePKLikeFun(idxvalue, 0, -1, true, cf.addtext)
	cf.Loop = 0                                                             //重置
	tdb.Rd[0] = xbdb.JoinBytes(tdb.Rd[0], []byte(xbdb.Split), cf.r.Bytes()) //kvs[0].V + xbdb.Split + cf.r.String() //在后面添加text的值
	if cf.r != nil {
		cf.r.Reset()
		bufpool.Put(cf.r)
	}
	tifo := Table[tbname].Ifo                                      //创建一个临时的表信息
	tifo.Fields = append(Table[tbname].Ifo.Fields, "text")         //添加一个text字段
	tifo.FieldType = append(Table[tbname].Ifo.FieldType, "string") //添加一个text字段类型
	r := Table[tbname].DataToJsonforIfo(tdb, &tifo)                //DataToJson(tdb, tifo)
	if r != nil {
		w.Write(r.Bytes())
		//w.Write([]byte(strconv.Quote(r.String()))) //必须使用strconv.Quote转义
		r.Reset()
		xbdb.Bufpool.Put(r)
	}

}

type cexefun struct {
	r      *bytes.Buffer
	tbname string
	bkey   []byte
	Loop   int
}

func newcexefun() *cexefun {
	tbname := "c"
	return &cexefun{
		r:      bufpool.Get().(*bytes.Buffer),
		tbname: tbname,
		bkey:   Table[tbname].Select.GetPkKeyLike([]byte("")), //获取前缀
	}
}
func (c *cexefun) addtext(rd []byte) bool {
	if c.Loop == 0 { //第一条是标题，省去。
		c.Loop++
		return true
	}
	vs := bytes.Split(rd, []byte(xbdb.Split))
	c.r.Write(vs[1])
	//c.r.Write([]byte("\n"))
	c.Loop++
	return true
}
