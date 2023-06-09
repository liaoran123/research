package routers

import (
	"bytes"
	"fmt"
	"net/http"
	"research/pubgo"
	"research/xbdb"
	"strconv"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Searchapp是为了blazor的GetFromJsonAsync设计的json格式。
// GetFromJsonAsync只能解析列表或嵌套列表
func SearchApp(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")
	pubgo.Tj.Brows("/SearchApp/")

	params := getparas(req)
	tbname := "c"
	params["kw"] = Sublen(params["kw"], 35) //最大长度35
	Se := NewSeExefunc(tbname, params["kw"], params["dir"], 21)
	asc := params["asc"] == "" //params["asc"]默认空值即true
	p := params["p"]
	ok := false

	var key []byte
	var iter iterator.Iterator

	if p == "" {
		//第一页搜索没有p值，需要查询获得
		key = Table[tbname].Select.GetIdxPrefixLike([]byte("s"), []byte(Se.mkw))
		iter, ok = Table[tbname].Select.IterPrefixMove(key, asc)
	} else { //整型转byte留下的复制问题
		//《一人心念口言--14371||0
		ps := strings.Split(p, xbdb.Split)
		if len(ps) > 1 {
			p = strings.Replace(p, "whfgf", "?", -1) //url转义
			p = strings.Replace(p, "yhfgf", "&", -1) //url转义
			ips := strings.Split(ps[1], idssplit)    //将14371+0转为byte的字符串
			if len(ips) > 1 {
				aid, _ := strconv.Atoi(ips[0])
				sid, _ := strconv.Atoi(ips[1])
				ids := ArtSecToId(aid, sid)
				//pos, _ := strconv.Atoi(ps[2])
				//key = JoinBytes([]byte(ps[0]), []byte(xbdb.Split), []byte(ids), []byte(xbdb.Split), IntToBytes(pos))
				key = JoinBytes([]byte(ps[0]), []byte(xbdb.Split), []byte(ids))
				key = Table[tbname].Select.GetIdxPrefixLike([]byte("s"), key)
				iter, ok = Table[tbname].Select.IterSeekMove([]byte(key))
			} else {
				ok = false
				fmt.Println("错误的定位页p,ps：", p, ps)
				fmt.Println("kw=" + params["kw"])
			}
		} else {
			ok = false
			fmt.Println("错误的定位页p：", p)
		}

	}
	if !ok {
		//fmt.Printf("p: %v\n", p)
		//fmt.Println("kw=" + params["kw"])
		return
	}
	//ts := pubgo.Newts() //计算执行时间
	Se.r.WriteString("[") //Se.r.WriteString("{\"datalist\":[")

	xbdb.NewIters(iter, ok, asc, 0, -1).ForKVFun(Se.searchapp)
	jsonstr := Se.r.String()
	jsonstr = strings.Trim(jsonstr, ",")
	Se.r.Reset()
	Se.r.WriteString(jsonstr)
	Se.r.WriteString("]") //Se.r.WriteString("],")
	/*
		Se.r.WriteString("\"lastkey\":" + strconv.Quote(Se.lastkey) + ",")
		setime := ts.Gstrts()
		Se.r.WriteString("\"setime\":\"" + setime + "\",")
		//fmt.Printf("setime: %v\n", setime)
		Se.r.WriteString("\"count\":" + strconv.Itoa(Se.loop) + "}")
	*/
	w.Write(Se.r.Bytes())
	//w.Write([]byte(strconv.Quote(Se.r.String()))) //必须使用strconv.Quote转义
	Se.r.Reset()
	bufpool.Put(Se.r)
}

// 搜索执行类
type SeExefunc struct {
	tbname  string
	kw      string
	dir     string   //目录范围，可以是多个
	ks      []string //用空格来判断组合查询，分解出多个词
	mkw     string   //最长的关键词
	count   int      //返回条数
	mlen    int      //摘录最大长度
	loop    int
	r       *bytes.Buffer
	lastkey string //最后的key值
	//--变量---
	//keys              [][]byte
	artid, secid, pos int
	cid               uint32
	bt                time.Time
}

func NewSeExefunc(tbname, kw, dir string, count int) *SeExefunc {
	ks := strings.Split(kw, " ") //用空格来判断组合查询
	//获取字数最长的词，通常字数最长的就是数据量最少的词。以该词作为组合查询的遍历定位词。
	mkw := getMaxLenKw(ks)
	//maxkeylen := int(ConfigMap["maxkeylen"].(float64))
	mkw = Sublen(mkw, 7) //搜索词最大长度是7
	return &SeExefunc{
		tbname: tbname,
		kw:     kw,
		dir:    dir,
		ks:     ks,
		mkw:    mkw,
		count:  count,
		mlen:   49,
		r:      bufpool.Get().(*bytes.Buffer),
		bt:     time.Now(),
	}
}
func (e *SeExefunc) searchapp(k, v []byte) bool {

	if time.Since(e.bt).Seconds() > 3 { //只要是控制组合查询超时时间
		e.loop = 21 //以便用户点击下一页，分散时间进行搜索，缓解性能问题。
		//fmt.Println("组合查询超时3秒。") //多次执行由于会加载内存，则可以完成。
		return false
	}

	if !strings.Contains(string(k), e.mkw) {
		return false //key不存在kw，即已经超过所需数据
	}
	//rd := xbdb.KVToRd(k, v, []int{})
	//解构rd，转为字符串lastkey。参照artpost.ArtSecToId组合规则
	keys := Table[e.tbname].Split(k) //bytes.Split(rd, []byte(xbdb.Split))
	artid, secid := IdToArtSec(string(keys[2]))
	if artid == 0 {
		return true
	}
	if (e.artid == artid) && (e.secid == secid) { //排除重复。同一段落包含多个相同kw时，出现重复情况。
		return true
	}
	e.artid = artid
	e.secid = secid
	e.pos = xbdb.BytesToInt(keys[2])

	e.lastkey = string(keys[1]) + xbdb.Split + strconv.Itoa(e.artid) + idssplit + strconv.Itoa(e.secid)
	e.lastkey = strings.Replace(e.lastkey, "?", "whfgf", -1) //url转义 ?和&不能出现再url。
	e.lastkey = strings.Replace(e.lastkey, "&", "yhfgf", -1) //url转义 ?和&不能出现再url。

	e.cid = Artfid[uint32(e.artid)]   //获取文章对应的所属目录id
	if !CacaRand(int(e.cid), e.dir) { //范围搜索
		return true
	}
	if !e.exsit() { //组合查询
		return true
	}
	e.r.WriteString("{\"dir\":" + CRAMs.GetCaDirToJsonApp(int(e.cid)) + ",") //写入目录路径
	e.r.WriteString(e.getartinfo() + ",")                                    //写入文章标题
	e.r.WriteString(e.getartmeta())                                          //写入文章摘录信息
	e.r.WriteString(",\"keyid\":" + strconv.Quote(e.lastkey) + "},")         //写入lastkey

	e.loop++
	return e.loop < e.count
}

// 组合查询
func (e *SeExefunc) exsit() (find bool) {
	if len(e.ks) < 2 {
		find = true
		return
	}
	id := ArtSecToId(e.artid, e.secid)
	idxvalue := Table[e.tbname].Ifo.FieldChByte("id", id)
	btext := Table[e.tbname].Select.GetPKValue(idxvalue)
	//sec := string(btext)
	secstr := string(btext)
	fc := 0
	for i := 0; i < len(e.ks); i++ { //for _, v := range e.ks { //如果在该段落内容里，所有的词组都存在，即是匹配。
		if strings.Contains(secstr, Sublen(e.ks[i], 7)) {
			fc++
		}
		//find = find && strings.Contains(secstr, Sublen(e.ks[i], 7)) //精准查询
	}
	if fc >= len(e.ks)/2+1 { //存在一半以上即当为匹配
		find = true
	}
	return
}
func (e *SeExefunc) getartinfo() (r string) {
	skid := ArtSecToId(e.artid, 0) //第0句是文章标题
	key := Table[e.tbname].Ifo.FieldChByte("id", skid)
	value := Table[e.tbname].Select.GetPKValue(key)
	/*
		[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
		{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
	*/
	r = "\"id\":" + strconv.Itoa(e.artid) + ","
	r += "\"title\":" + strconv.Quote(string(value)) //+ "\""

	return
}

// 文章摘录
func (e *SeExefunc) getartmeta() (r string) {
	skid := ArtSecToId(e.artid, e.secid)                                 //c表id的字符串
	key := Table[e.tbname].Select.GetPkKey(xbdb.SplitToCh([]byte(skid))) //Table[e.tbname].Ifo.FieldChByte("id", skid)
	iter, ok := Table[e.tbname].Select.IterSeekMove(key)
	if !ok {
		return
	}
	meta := string(iter.Value())
	eid := e.secid
	for len(meta) < e.mlen*3 { //每个中文3个字节
		eid++
		ok = iter.Next()
		if !ok {
			break
		}
		meta += string(iter.Value())
	}

	/*
		[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
		{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
	*/
	r = "\"bid\":" + strconv.Itoa(e.secid) + ","
	r += "\"eid\":" + strconv.Itoa(eid) + ","
	r += "\"text\":" + strconv.Quote(meta) //+ "\""

	return
}

/*
// 写超时警告日志 通用方法

func TimeoutWarning(tag, detailed string, start time.Time, timeLimit float64) {
	dis := time.Now().Sub(start).Seconds()
	if dis > timeLimit {
		log.Warning(log.CENTER_COMMON_WARNING, tag, " detailed:", detailed, "TimeoutWarning using", dis, "s")
		//pubstr := fmt.Sprintf("%s count %v, using %f seconds", tag, count, dis)
		//stats.Publish(tag, pubstr)
	}
}
*/
//找出最大长度的词
func getMaxLenKw(ks []string) (s string) {
	l := 0
	lv := 0
	for _, v := range ks {
		lv = len([]rune(v))
		if lv >= l {
			s = v
			l = lv
		}
	}
	return
}
