package base

import (
	"research/pubgo"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var (
	fixed map[bool]func(iter iterator.Iterator) bool //起始位置。升序，first；降序，last
	move  map[bool]func(iter iterator.Iterator) bool //遍历顺序动作。升序，next；降序，prev
)
var PSe *Se

//搜索结构体
type Se struct {
	cont *content
	//order bool //ture升序，反之。
	//pgcount int
}

func NewSe() *Se {
	if fixed == nil {
		fixed = make(map[bool]func(iter iterator.Iterator) bool, 2)
		fixed[true] = first
		fixed[false] = last
	}
	if move == nil {
		move = make(map[bool]func(iter iterator.Iterator) bool, 2)
		move[true] = next
		move[false] = prev
	}

	return &Se{
		cont: Newcontent(),
		//order: order,
		//pgcount: pgcount,
	}
}

//搜索

func (s *Se) Search(keyword, p, count, caids string, order bool) (rsetAll *RsetAll) {
	pgcount, _ := strconv.Atoi(count)
	if pgcount == 0 {
		pgcount = 21 //默认21条记录
	}
	ids, lastkey, SeTime := s.cont.Search(keyword, p, caids, order, pgcount)
	ts := pubgo.Newts() //计算执行时间
	/*
		只需要实现 New 函数,
		Get()时，对象池中没有对象，将会自动调用 New 函数创建。
	*/
	minlentext := 21
	if ConfigMap["minlentext"] != nil {
		minlentext = int(ConfigMap["minlentext"].(float64))
	}
	rsetAll = RsetAllPool.Get().(*RsetAll)
	rset := RsetPool.Get().(*Rset)
	//var r *cataRAM
	rsetAll.SeTime = SeTime
	for _, iv := range ids {
		//rset.Cataid = iv.cataid
		rset.CataDir = CRAMs.GetCataDir(iv.cataid)
		/*
			if rset.CataDir == nil {//用户可能没有设置目录
				fmt.Println("0")
			}*/

		rset.Artid = iv.artid
		if iv.artid == 0 { //用户未将分隔符"-"转义时会有这样的情况。
			continue
		}
		rset.Secid = iv.secid
		rset.Title, rset.ArtUrl, rset.Text, rset.LastSecid = s.cont.GetArtPathInfo(iv.artid, iv.secid, minlentext)
		rset.Text = strings.Replace(rset.Text, "﹣", "-", -1) //text = strings.Replace(text, "-", "﹣", -1) //-是系统保留字，需要转义为﹣。
		if rset.ArtUrl == "0" {
			rset.ArtUrl = ""
		}
		rsetAll.Set = append(rsetAll.Set, *rset)
	}
	RsetPool.Put(rset)
	rsetAll.Lastkey = lastkey
	SetTime := ts.Gstrts()
	rsetAll.SetTime = SetTime
	//fmt.Println("Se:" + ys)
	return
}

/*
//检测目录有效性。
func (s *Se) chcata(caids string) (r bool) {
	ids := strings.Split(caids, "|")
	iv := 0
	for _, v := range ids {
		iv, _ = strconv.Atoi(v)
		if _, r = CRAMs.CataRAMMap[uint32(iv)]; r {
			return
		}
	}
	return
}


//获取指定最少minlentext字数的连续内容作为搜索返回结果节录。
func (s *Se) GetMinLenText(iter iterator.Iterator, artid, secid int) (r string, LastSecid int) {
	minlentext := 21
	if ConfigMap["minlentext"] != nil {
		minlentext = int(ConfigMap["minlentext"].(float64))
	}
	r, LastSecid = s.cont.IterMinLenText(iter, artid, secid, minlentext)
	return
}


//获取指定最少minlentext字数的连续内容作为搜索返回结果节录。
func (s *Se) GetMinLenText(artid, secid int) (r string, LastSecid int) {
	minlentext := 21
	if ConfigMap["minlentext"] != nil {
		minlentext = int(ConfigMap["minlentext"].(float64))
	}
	r, LastSecid = s.cont.GetMinLenText(artid, secid, minlentext)
	return
}*/
func first(iter iterator.Iterator) bool {
	//fmt.Println("First")
	return iter.First()
}
func last(iter iterator.Iterator) bool {
	//fmt.Println("Last")
	return iter.Last()
}
func prev(iter iterator.Iterator) bool {
	//fmt.Println("Prev")
	return iter.Prev()
}
func next(iter iterator.Iterator) bool {
	//fmt.Println("Next")
	return iter.Next()
}
