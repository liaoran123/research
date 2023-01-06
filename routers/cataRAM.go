package routers

import (
	"bytes"
	"research/xbdb"
	"strconv"
	"strings"
)

var CRAMs *CataRAMs

//目录数据少，使用频繁，故适宜加载入内存
//删除目录时不能真正删除，只需将除id外置空即可。这样保证通过id直接匹配数组下标。
type cataRAM struct {
	/*id,*/ fid   int
	title, isleaf string
}

func NewcataRAM(fid int, title, isleaf string) *cataRAM {
	return &cataRAM{
		//id:   id,
		fid:    fid,
		title:  title,
		isleaf: isleaf,
	}
}

type CataRAMs struct {
	CataRAMMap map[uint32]uint32 //记录id对应的cataRAM数组下标。目录id并不是顺序的。同时可以减少gc
	cataRAM    []*cataRAM
}

func NewCataRAMs() *CataRAMs {
	return &CataRAMs{
		CataRAMMap: make(map[uint32]uint32),
	}
}
func (c *CataRAMs) LoadCataRAM() {
	key := []byte("ca" + xbdb.Split)
	tbd := Table["ca"].Select.FindPrefix(key, true, 0, -1)
	if tbd != nil {
		c.toRAM(tbd)
	}

}
func (c *CataRAMs) toRAM(tbd *xbdb.TbData) {
	for _, v := range tbd.Rd {
		key := bytes.Split(v, []byte(xbdb.Split))
		id := BytesToInt(key[0])
		title := key[1]
		fid := BytesToInt(key[2])
		isleaf := key[3]
		c.Append(id, fid, string(title), string(isleaf))
	}
}
func (c *CataRAMs) Append(id, fid int, title, isleaf string) {
	if v, ok := c.CataRAMMap[uint32(id)]; !ok {
		c.cataRAM = append(c.cataRAM, NewcataRAM(fid, title, isleaf))
		c.CataRAMMap[uint32(id)] = uint32(len(c.cataRAM)) - 1 //记录id对应的数组下标
	} else { //如果存在即修改。事务回滚会有出现这个情况。
		c.cataRAM[v].fid = fid
		c.cataRAM[v].title = title
	}
}

//通过目录id获取一个目录信息
func (c *CataRAMs) Get(id int) (r *cataRAM) {
	if id, ok := c.CataRAMMap[uint32(id)]; ok {
		r = c.cataRAM[id]
	}
	return
}

/*
//获取目录路径
func (c *CataRAMs) GetCataDir(cataid int) (r []CataInfo) {
	cd := CataInfo{}
	var ok bool
	var i uint32
	cid := cataid
	L := 0
	for cid > 0 {
		if i, ok = CRAMs.CataRAMMap[uint32(cid)]; ok {
			cd.Id = cid
			cd.Title = CRAMs.cataRAM[i].title
			cd.Fid = CRAMs.cataRAM[i].fid
			cd.Isleaf = CRAMs.cataRAM[i].isleaf
			cid = cd.Fid //CRAMs.cataRAM[i].fid
			r = append(r, cd)
		} else { //用户没有设置目录
			return
		}
		L++ //以防用户目录混乱导致的死循环
		if L > 49 {
			return
		}
	}
	return
}
*/
//获取目录路径
func (c *CataRAMs) GetCaDirToJson(cataid int) (r string) {
	var ok bool
	var i uint32
	cid := cataid
	L := 0
	jsonstr := "{\"result\":["
	for cid > 0 {
		if i, ok = CRAMs.CataRAMMap[uint32(cid)]; ok {
			/*
				[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
				{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
			*/
			jsonstr += "{\"id\":" + strconv.Itoa(cid) + ","
			jsonstr += "\"title\":" + strconv.Quote(CRAMs.cataRAM[i].title) + ","
			jsonstr += "\"fid\":" + strconv.Itoa(CRAMs.cataRAM[i].fid) + ","
			jsonstr += "\"isleaf\":\"" + CRAMs.cataRAM[i].isleaf + "\"},"
			cid = CRAMs.cataRAM[i].fid

		} else { //用户没有设置目录
			break
		}
		L++ //以防用户目录混乱导致的死循环
		if L > 49 {
			break
		}
	}
	r = strings.Trim(jsonstr, ",") + "]}"
	return
}

//在某个或多个目录下查找
//caids目录id集合
func CacaRand(caid int, caids string) (r bool) {
	if caids == "" || caid == 0 {
		r = true
		return
	}
	ids := "|" + caids + "|"

	fid := caid //CRAMs.Get(artid - 1).fid //CRAMs.cataRAM[artid-1].fid
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

/*

//并不是真正删除。否则目录表id和数组下标不能对应。
func (c *CataRAMs) Del(id int) {
	if id, ok := c.CataRAMMap[uint32(id)]; ok {
		c.cataRAM[id].fid = 0
		c.cataRAM[id].title = ""
		c.cataRAM[id].isleaf = ""
	}
}

//通过目录id修改。
func (c *CataRAMs) Put(id, fid int, title, isleaf string) {
	if id, ok := c.CataRAMMap[uint32(id)]; ok {
		c.cataRAM[id].fid = fid
		c.cataRAM[id].title = title
	}
}
func (c *CataRAMs) InsOrUpd(params map[string]string) {
	//id, fid int, title, isleaf string
	id, _ := strconv.Atoi(params["id"])
	fid, _ := strconv.Atoi(params["fid"])
	title := params["title"]
	isleaf := params["isleaf"]
	if params["iou"] == "ins" {
		c.Append(id, fid, title, isleaf)
	} else {
		c.Put(id, fid, title, isleaf)
	}
}
*/
