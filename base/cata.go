package base

import (
	"strings"
	"sync"
)

//每个表对应一个结构体
var Pcata *cata

//目录表
type cata struct {
	tbn   string
	id    int
	cfidx *catafididx
	ainc  *Autoinc
	mu    sync.RWMutex
}

func Newcata() *cata {
	if PcataAutoinc == nil {
		PcataAutoinc = NewAutoinc("ca")
	}
	return &cata{
		tbn:   "ca", //表名
		cfidx: Newcatafididx(),
		ainc:  PcataAutoinc, //NewAutoinc("ca"),
	}
}

//目录数据少，使用频繁，故适宜加载入内存
//删除目录时不能真正删除，只需将除id外置空即可。这样保证通过id直接匹配数组下标。

//添加时用
func (c *cata) GetAutotid() {
	c.id = c.ainc.Getid()
}

//修改删除用
func (c *cata) Setid(id int) {
	c.id = id
}

//添加目录
func (c *cata) Insert(id int, title, isleaf string, fid int) (r bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	r = c.Insertcata(id, title, isleaf, fid) //经过这里id=c.id
	r = r && c.cfidx.Insert(fid, c.id)
	//同时需要添加一篇只有目录标题而没有内容的文章，以便通过内容来搜索到目录。
	r = r && PArticle.Insert(0, title, "", "。", "0", id) //文章的fid就是目录的id。Insert(0...，0,表示id是自动增值。
	if r {                                               //如果成功，添加目录到内存
		CRAMs.Append(c.id, fid, title, "0") //= append(CRAMs.cataRAM, NewcataRAM(c.id, fid, title)) //实时加入内存
		//c.ainc.Writelastid()
	}
	return
}

//添加目录
func (c *cata) Insertcata(id int, name, isleaf string, fid int) (r bool) {
	if id == 0 { //不传id则自动增值。因为需要批量导入，id由客户决定。
		c.GetAutotid()
	} else {
		c.id = id
	}
	/*
			//k=ca-id
			//v=name-fid
			//setkey

		//fmt.Println(id, name, isleaf, fid)
		if id == 42 {
			fmt.Print("fff")
		}
	*/
	err = Con.Getartdb().Db.Put(c.setkey(c.id), JoinBytes([]byte(name+"~"), IntToBytes(fid), []byte("~"+isleaf)), nil) //添加目录标题
	//err = Con.Getartdb().Db.Put(c.setkey(c.id), []byte(name+"~"+fmt.Sprintf("%11s", strconv.Itoa(fid))+"~"+isleaf), nil) //添加目录标题
	Chekerr()
	r = err == nil
	return
}

func (c *cata) Deletetcata(id int) (r bool) {
	err = Con.Getartdb().Db.Delete(c.setkey(id), nil)
	Chekerr()
	r = err == nil
	return
}

//-----------查询-------------------
//获取一个目录的信息
func (c *cata) GetCata(id int) (r CataInfo) {
	data, _ := Con.Getartdb().Db.Get(c.setkey(id), nil)
	if data == nil {
		return
	}
	sdata := strings.Split(string(data), "~")
	r.Id = id
	r.Name = sdata[0]
	if r.Name == "" {
		return
	}
	r.Fid = BytesToInt([]byte(sdata[1]))
	r.Isleaf = sdata[2]
	return
}

//返回所有子目录信息
func (c *cata) ChildCatas(fid int) (r []CataInfo) {
	ids := c.cfidx.ChildCatas(fid)
	for _, iv := range ids { //实时添加到内存
		r = append(r, c.GetCata(iv))
	}
	return
}
func (c *cata) setkey(id int) (r []byte) {
	r = JoinBytes([]byte(c.tbn+"~"), IntToBytes(id)) //[]byte(c.tbn + "~" + fmt.Sprintf("%11s", strconv.Itoa(id))) //
	return
}

//删除目录
func (c *cata) Delete(id int) (r bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	r = c.Deletetcata(id)
	fid := CRAMs.Get(id).fid
	r = r && c.cfidx.Delete(fid, id)
	//同时需要添加一篇只有目录标题而没有内容的文章，以便通过内容来搜索到目录。
	r = r && NewArticle().Delete(id)

	if r { //如果成功，添加目录到内存
		CRAMs.Del(id)
	}
	return
}

//修改目录
func (c *cata) Put(id, fid int, name string) (r bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k := c.setkey(id) //JoinBytes([]byte(c.tbn+"~"), IntToBytes(id))
	v := JoinBytes([]byte(name+"~"), IntToBytes(fid), []byte("~"))
	err = Con.Getartdb().Db.Put(k, v, nil)
	Chekerr()
	r = err == nil
	if r { //修改内存
		CRAMs.Put(id, fid, name)
	}
	return
}
func (c *cata) GetName(id int) (name string, fid int) {
	data, _ := Con.Getartdb().Db.Get(c.setkey(id), nil)
	if data == nil { //顶级目录没有数据
		return
	}
	ds := strings.Split(string(data), "~")
	name = ds[0]
	sfid := ds[1]
	fid = BytesToInt([]byte(sfid))
	return
}
