package base

/*
i- 内容索引表
.k=i-长度7的分段内容遍历分词-文章id-分段id;v=空值
.i-最后后面加内容表的唯一id："文章id-分段id"，为了相同分词之后按内容表先后排序一致
*/
import (
	"bytes"
	"research/pubgo"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/util"
)

var Pcontent *content

//内容表
type content struct {
	id      int
	tbn     string
	Idx     *idx
	fidx    *fididx
	split   string //组合查询分隔符
	maxloop int    //慢查询规避。结果皆不匹配的情况下，组合查询最大匹配多少次。
	ainc    *Autoinc
}

func Newcontent() *content {
	split := " " //组合查询分隔符。默认是空格
	if ConfigMap["seslipt"] != nil {
		split = ConfigMap["seslipt"].(string)
	}
	maxloop := 0 //默认不设置组合慢查询规避
	if ConfigMap["maxloop"] != nil {
		maxloop = int(ConfigMap["maxloop"].(float64))
	}
	if PcontenAutoinc == nil {
		PcontenAutoinc = NewAutoinc("c")
	}
	return &content{
		tbn:     "c",
		Idx:     Newidx(),
		fidx:    Newfididx(),
		split:   split,
		maxloop: maxloop,
		ainc:    PcontenAutoinc, //NewAutoinc("c"),
	}
}

//***********添加*******************
//添加内容，将文章分成多个句子段落后添加到表
func (c *content) TextSplit(text, split string) (section []string) {
	itext := text //title+"\n"+text
	if split != "" {
		//支持多个分段匹配标签。中文常见是“。”.
		//空格是组合查询，由于支持英文，故而空格不作默认分隔符
		ss := strings.Split(split, "|")
		for _, v := range ss {
			itext = strings.Replace(itext, v, v+"\n", -1) //分配段落
		}
	}
	section = strings.Split(itext, "\n")
	return
}

//添加内容
//id支持自动增值或用户定义。
//用户定义主要是用于批量导入，可以保证与原数据保持一致。
//id=32位无符号整型。uint32
func (c *content) Insert(id, fcataid int, title, text, split, url string) (r bool) {
	if id == 0 { //自动增值默认从1开始。传入的id=0，即表示使用自动增值。
		c.id = c.ainc.Getid()
	} else {
		c.id = id
	}

	r = c.InsertFid(c.id, fcataid)
	r = r && c.InsertConAndIdx(c.id, title, text, split, url, fcataid)
	/*
		if r {
			c.ainc.Writelastid()
		}*/
	return
}

//用于修改的添加
func (c *content) InsertId(id, fid int, title, text, split, url string) (r bool) {
	r = c.InsertFid(id, fid)
	r = r && c.InsertConAndIdx(id, title, text, split, url, fid)
	return
}

////打开文章遍历段落，删除索引和内容
func (c *content) Delete(id int) (r bool) {
	fid := c.GetArtFCataId(id) //ReFid[id] //文章id找fid即目录id
	r = c.DeleteFid(id, fid)
	r = r && c.DeleteConAndIdx(id)
	return

}

//添加文章对应的fid，即所属目录
func (c *content) InsertFid(artid, fid int) (r bool) {
	/*
		if fid != 0 { //0，是顶级目录
			r = c.fidx.Insert(fid, artid)
		} else {
			r = true
		}*/
	r = c.fidx.Insert(fid, artid)
	return
}

//添加文章对应的fid，即所属目录
func (c *content) DeleteFid(artid, fid int) (r bool) {
	/*
		if fid != 0 { //0，是顶级目录
			r = c.fidx.Delete(fid, artid)
		} else {
			r = true
		}*/
	r = c.fidx.Delete(fid, artid)
	return
}

//添加内容和索引，将文章分成多个句子段落后添加到表和索引
func (c *content) InsertConAndIdx(id int, title, text, split, url string, fcataid int) (r bool) {
	i := 0
	s := ""
	u := "0" //当url为空时，以此为标志。
	if url != "" {
		u = url
	}
	/*
		《"+title+"》\n 将标题加入内容，即不需要多加一个标题的搜索。通过“《”+关键词就能专门搜索标题
		文章内容约定，0段落=标题；1段落分隔符=split；2段落=url；3段落=fcataid。
	*/
	section := strings.Split(("《" + title + "》\n" + split + "\n" + u + "\n" + strconv.Itoa(fcataid)), "\n")
	section = append(section, c.TextSplit(text, split)...)
	r = true
	for _, sec := range section {
		s = strings.TrimSpace(sec)
		if s == "" {
			continue
		}
		if r { //添加内容
			//err = Con.Getartdb().Db.Put(JoinBytes([]byte(c.tbn+"-"), IntToBytes(id), []byte("-"), IntToBytes(i)), []byte(s), nil)
			err = Con.Getartdb().Db.Put(c.setkey(id, i), []byte(s), nil)
		} else {
			return
		}
		r = r && err == nil
		r = r && c.Idx.Act(id, i, fcataid, sec, c.Idx.Insert) //添加内容段落索引
		i++
		Chekerr()
	}
	return
}

//key=c-artid-secid
func (c *content) setkey(artid, secid int) (r []byte) {
	r = JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid))
	return
}

//删除内容
func (c *content) DeleteConAndIdx(id int) (r bool) {
	r = true
	Prefix := JoinBytes([]byte(c.tbn+"-"), IntToBytes(id), []byte("-"))
	iter := Con.Getartdb().Db.NewIterator(util.BytesPrefix(Prefix), nil) //打开文章所有段落，遍历删除
	var ks []string
	var artid, secid int
	for iter.Next() {
		ks = strings.Split(string(iter.Key()), "-")
		artid = BytesToInt([]byte(ks[1]))
		secid = BytesToInt([]byte(ks[2]))
		r = c.Idx.Act(artid, secid, 0, string(iter.Value()), c.Idx.Delete) //删除索引fcataid int 参数=0。删除只需key，不需要value。与fcataid无关。
		if r {
			err = Con.Getartdb().Db.Delete(iter.Key(), nil) //删除段落
		}
		Chekerr()
	}
	Release(iter)
	return
}

//获取文章的所属目录id
func (c *content) GetArtFCataId(id int) (r int) {
	//文章内容约定，0段落=标题；1段落分隔符=split；2段落=url；3段落=fcataid。
	Prefix := JoinBytes([]byte(c.tbn+"-"), IntToBytes(id), []byte("-"), IntToBytes(3))
	data, _ := Con.Getartdb().Db.Get(Prefix, nil)
	r, _ = strconv.Atoi(string(data))
	return
}

//获取一文章
func (c *content) GetArtInfo(id int) (r Artinfo) {
	Prefix := JoinBytes([]byte(c.tbn+"-"), IntToBytes(id), []byte("-"))
	iter := Con.Getartdb().Db.NewIterator(util.BytesPrefix(Prefix), nil)
	i := -1
	var title, split, url, text string
	var fid int
	buf := ArtBuf.Get().(*bytes.Buffer)
	//文章内容约定，0段落=标题；1段落分隔符=split；2段落=url；3段落=fcataid。
	for iter.Next() {
		i++
		if i == 0 {
			title = string(iter.Value())
			title = strings.Replace(title, "《", "", -1)
			title = strings.Replace(title, "》", "", -1)
			continue
		}
		if i == 1 {
			split = string(iter.Value())
			continue
		}
		if i == 2 {
			url = string(iter.Value())
			if url == "0" {
				url = ""
			}
			continue
		}
		if i == 3 {
			sfid := string(iter.Value())
			fid, _ = strconv.Atoi(sfid)
			continue
		}
		buf.Write(iter.Value())
	}
	Release(iter)
	text = buf.String()
	buf.Reset()
	ArtBuf.Put(buf)
	//文章内容约定，0段落=标题；1段落分隔符=split；2段落=url；3段落=fcataid。
	r.Title = title
	r.Split = split
	r.Fid = fid
	r.Url = url
	r.Text = strings.Replace(text, "﹣", "-", -1) //text = strings.Replace(text, "-", "﹣", -1) //-是系统保留字，需要转义为﹣。
	return
}

//******搜索*************************************
//idx.Search(...)该函数不能放进idx结构类。
//为了优化组合查询速度，搜索时会使用到content类。
//这样造成循环引用。故提出到该层来实现。

func (c *content) Search(keyword, p, caids string, order bool, count int) (asids []asid, lastkey, time string) {
	ts := pubgo.Newts() //计算执行时间
	ks := strings.Split(keyword, c.split)
	//获取字数最长的词，通常字数最长的就是数据量最少的词。以该词作为组合查询的遍历定位词。
	kw := c.getMaxLenKw(ks)
	kw = Sublen(kw, c.Idx.keywordlen) //搜索词最大长度
	var jl int                        //记录组合查询的次数，避免大数据库而无符合条件的情况。
	var je bool
	var ok bool
	iter := c.Idx.Getiter(kw) //获取kw为前缀的查询索引游标。
	if p != "" {
		//重新组成key。
		//golang json.Marshal 特殊html字符被转义
		r := c.Joinkey(p)
		if r != nil {
			ok = iter.Seek(r) //第二页开始可以定位。
		} else { //p错误的时候，应该是认为问题.
			return
		}
	} else {
		ok = fixed[order](iter) //第一页定位。升序first，降序last。
	}

	succ := 0
	var e bool
	var key string
	var keys []string
	var cataid, artid, secid int
	tasid := asid{}
	for ok {
		e = true
		key = string(iter.Key())
		// - 该符号需要转义，最好使用\n分隔，不会冲突。
		keys = strings.Split(key, "-")
		//ckw = keys[1]
		artid = BytesToInt([]byte(keys[2])) //文章id
		secid = BytesToInt([]byte(keys[3])) //段落id
		lastkey = c.setlastkey(keys[1], artid, secid)
		cataid = BytesToInt(iter.Value()) //iter.Value()就是文章的目录id
		e = e && c.rand(cataid, caids)    //目录范围查询
		je, jl = c.exsit(artid, secid, ks, jl)
		if c.maxloop != 0 && jl > c.maxloop { //组合查询最多1024次，没有符合条件即退出。避免慢查询。
			//fmt.Println(keyword, 1024)
			return
		}
		e = e && je //c.exsit(artid, secid, ks, jl) //组合查询
		if !e {
			ok = move[order](iter)
			continue
		}
		if (tasid.artid == artid) && (tasid.secid == secid) {
			ok = move[order](iter)
			continue //排除重复。同一段落包含多个相同kw时，出现重复情况。
		}
		tasid.cataid = cataid
		tasid.artid = artid
		tasid.secid = secid
		asids = append(asids, tasid)
		succ++
		if succ >= count {
			break
		}
		ok = move[order](iter) //游标移动。升序next，降序prev。
	}
	Release(iter)
	time = ts.Gstrts()
	return
}

//将lastkey转换为可视化字符串。
func (c *content) setlastkey(kw string, artid, secid int) string {
	a := strconv.Itoa(artid)
	s := strconv.Itoa(secid)
	return kw + "," + a + "," + s
}

//在某个或多个目录下查找
//caids目录id集合
func (c *content) rand(cataid int, caids string) (r bool) {
	if caids == "" || cataid == 0 {
		r = true
		return
	}
	ids := "|" + caids + "|"

	fid := cataid //CRAMs.Get(artid - 1).fid //CRAMs.cataRAM[artid-1].fid
	loop := 0
	for fid > 0 { //遍历到顶级目录
		if strings.Contains(ids, "|"+strconv.Itoa(fid)+"|") {
			r = true
			return
		} else {
			if v, ok := CRAMs.CataRAMMap[uint32(fid)]; !ok { //防止用户输入的目录混乱。
				return
			} else {
				fid = CRAMs.cataRAM[v].fid
			}
		}
		loop++
		if loop >= 108 { //防止用户输入的目录混乱导致死循环。
			return
		}
	}
	return
}

//组合查询
func (c *content) exsit(artid, secid int, ks []string, lp int) (find bool, loop int) {
	find = true
	if len(ks) < 2 {
		return
	}
	sec := c.GetOneSec(artid, secid)
	loop = lp + 1
	Chekerr()
	secstr := string(sec)
	for _, v := range ks { //如果在该段落内容里，所有的词组都存在，即是匹配。
		//Sublen(v, c.idx.keywordlen),只需前面7个字（7=索引长度keywordlen）匹配即可。
		find = find && strings.Contains(secstr, Sublen(v, c.Idx.keywordlen))
	}
	return
}

//找出最大长度的词
func (c *content) getMaxLenKw(ks []string) (s string) {
	l := 0
	lv := 0
	for _, v := range ks {
		lv = len([]rune(v))
		if lv > l {
			s = v
			l = lv
		}
	}
	return
}
func (c *content) Joinkey(p string) (r []byte) {
	//重新组成key。
	//golang json.Marshal 特殊html字符被转义
	tp := strings.Replace(p, "\\u003c", "<", -1)
	tp = strings.Replace(tp, "\\u003e", ">", -1)
	tp = strings.Replace(tp, "\\u0026", "&", -1)
	ps := strings.Split(tp, ",")
	if len(ps) > 2 {
		aid, _ := strconv.Atoi(ps[1])
		sid, _ := strconv.Atoi(ps[2])
		r = JoinBytes([]byte("i-"), []byte(ps[0]), []byte("-"), IntToBytes(aid), []byte("-"), IntToBytes(sid))
	}
	return
}

/*
//byte的key中转换提取文章id，段落id
func (c *content) getasid(keys []string) (artid, secid int) {
	artid = BytesToInt([]byte(keys[2])) //文章id
	secid = BytesToInt([]byte(keys[3])) //段落id
	return
}*/

//通过文章id和段落id，获取段落id的内容。
func (c *content) GetOneSec(artid, secid int) (sec []byte) {
	if artid == 0 {
		return
	}
	sec, err = Con.Getartdb().Db.Get(JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid)), nil)
	Chekerr()
	return
}

/*
//通过文章id和段落id，获取段落id的内容。
func (c *content) IterOneSec(iter iterator.Iterator, artid, secid int) (sec []byte) {
	if artid == 0 {
		return
	}
	//sec, err = Con.Getartdb().Db.Get(JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid)), nil)
	k := JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid))
	if iter.Seek(k) {
		sec = iter.Value()
	}
	Chekerr()
	return
}*/

//获取搜索结果的节录内容
//在当前段落id至向后10个段落区间读取内容累加。直至内容长度达到最小长度minlentext。
func (c *content) GetArtPathInfo(artid, secid, minlentext int) (title, url, text string, LastSecid int) {

	//ts := pubgo.Newts() //计算执行时间
	iter := Con.Getartdb().Db.NewIterator(nil, nil)
	k := JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(0))
	ok := iter.Seek(k)
	if ok {
		title = string(iter.Value()) //文章标题
	}
	k = JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(2))
	ok = iter.Seek(k)
	if ok {
		url = string(iter.Value()) //文章网址
	}
	k = JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid))
	ok = iter.Seek(k)
	if !ok {
		return
	}
	var value string
	var ks []string
	var iaid, isid int //
	for ok {
		ks = strings.Split(string(iter.Key()), "-")
		if len(ks) < 3 {
			return
		}
		iaid = BytesToInt([]byte(ks[1]))
		isid = BytesToInt([]byte(ks[2]))
		//fmt.Println(iaid, isid)
		if iaid != artid { //已经是不同文章
			break
		}
		value = string(iter.Value())
		text += value + "\n"
		if len([]rune(text)) >= minlentext {
			break
		}
		iter.Next()
	}
	Release(iter)
	LastSecid = isid
	//ys := ts.Gstrts()
	//fmt.Println(artid, secid, ys)
	return
}

/*
//获取搜索结果的节录内容
//在当前段落id至向后10个段落区间读取内容累加。直至内容长度达到最小长度minlentext。
func (c *content) IterMinLenText(iter iterator.Iterator, artid, secid, minlentext int) (r string, LastSecid int) {

	k := JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid))
	ok := iter.Seek(k)
	if !ok {
		return
	}
	var value string
	var ks []string
	var iaid, isid int //
	for ok {
		ks = strings.Split(string(iter.Key()), "-")
		iaid = BytesToInt([]byte(ks[1]))
		isid = BytesToInt([]byte(ks[2]))
		//fmt.Println(iaid, isid)
		if iaid != artid { //已经是不同文章
			break
		}
		value = string(iter.Value())
		r += value + "\n"
		if len([]rune(r)) >= minlentext {
			break
		}
		iter.Next()
	}
	Release(iter)
	LastSecid = isid
	//ys := ts.Gstrts()
	//fmt.Println(artid, secid, ys)
	return
}

//返回一个空游标
func (c *content) Getartniliter() iterator.Iterator {
	return Con.Getartdb().Db.NewIterator(nil, nil)
}
*/
/*
func (c *content) GetMinLenText(artid, secid, minlentext int) (r string, LastSecid int) {
	b := JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid))
	e := JoinBytes([]byte(c.tbn+"-"), IntToBytes(artid), []byte("-"), IntToBytes(secid+10)) //最多10句，不够则忽略。
	iter := Con.Getartdb().Db.NewIterator(&util.Range{Start: []byte(b), Limit: []byte(e)}, nil)
	var value string
	var ks []string
	var iaid, isid int //
	for iter.Next() {
		ks = strings.Split(string(iter.Key()), "-")
		iaid = BytesToInt([]byte(ks[1]))
		isid = BytesToInt([]byte(ks[2]))
		if iaid != artid { //已经是不同文章
			break
		}
		value = string(iter.Value())

			r += value + "\n"
			if len([]rune(r)) >= minlentext {
				break
			}
		}
		Release(iter)
		LastSecid = isid
		return
	}

*/
