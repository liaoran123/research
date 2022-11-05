package base

import (
	"strings"
)

var CRAMs *CataRAMs

//目录数据少，使用频繁，故适宜加载入内存
//删除目录时不能真正删除，只需将除id外置空即可。这样保证通过id直接匹配数组下标。
type cataRAM struct {
	/*id,*/ fid int
	name        string
}

func NewcataRAM(fid int, name string) *cataRAM {
	return &cataRAM{
		//id:   id,
		fid:  fid,
		name: name,
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
	Con.Getartdb().FindRangeFun("ca-", "ca-a", c.toRAM)
}
func (c *CataRAMs) toRAM(k, v []byte) {
	key := strings.Split(string(k), "-")
	id := BytesToInt([]byte(key[1]))
	value := strings.Split(string(v), "-")

	fid := BytesToInt([]byte(value[1]))
	name := value[0]

	c.Append(id, fid, name)
	//c.cataRAM = append(c.cataRAM, NewcataRAM(fid, name))
	//c.CataRAMMap[uint32(id)] = uint32(len(c.cataRAM)) //记录id对应的数组下标
}
func (c *CataRAMs) Append(id, fid int, name string) {
	if v, ok := c.CataRAMMap[uint32(id)]; !ok {
		c.cataRAM = append(c.cataRAM, NewcataRAM(fid, name))
		c.CataRAMMap[uint32(id)] = uint32(len(c.cataRAM)) - 1 //记录id对应的数组下标
	} else { //如果存在即修改。事务回滚会有出现这个情况。
		c.cataRAM[v].fid = fid
		c.cataRAM[v].name = name
	}
}

//通过目录id获取一个目录信息
func (c *CataRAMs) Get(id int) (r *cataRAM) {
	if id, ok := c.CataRAMMap[uint32(id)]; ok {
		r = c.cataRAM[id]
	}
	return
}

//并不是真正删除。否则目录表id和数组下标不能对应。
func (c *CataRAMs) Del(id int) {
	if id, ok := c.CataRAMMap[uint32(id)]; ok {
		c.cataRAM[id].fid = 0
		c.cataRAM[id].name = ""
	}
}

//通过目录id修改。
func (c *CataRAMs) Put(id, fid int, name string) {
	if id, ok := c.CataRAMMap[uint32(id)]; ok {
		c.cataRAM[id].fid = fid
		c.cataRAM[id].name = name
	}
}

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
			cd.Name = CRAMs.cataRAM[i].name
			cd.Fid = CRAMs.cataRAM[i].fid
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
