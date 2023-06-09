package xbdb

import (
	"bytes"
	"fmt"
	"strconv"
)

// 数据库备份，全库或单表等导出备份
type Dbback struct {
	Oxb, Nxb           *Xb               //原新数据库和新数据库
	otbnames, ntbnames []string          //导入导出的所有表名
	otables, ntables   map[string]*Table //导入导出的所有表操作结构
	//以下3个参数是为了PutData作为参数而设。
	//每执行一次都要初始化一次这3参数。
	Upfield          map[string]string //要修改的字段名
	Otbname, Ntbname string            //导入导出的表名
}

// 根据路径获取数据库连接，用于独立使用，比如数据库管理桌面。
func Newdbback(odbpath, ndbpath string) *Dbback {
	oxb := NewDb(odbpath)
	nxb := NewDb(ndbpath)
	return &Dbback{
		Oxb:      oxb,
		Nxb:      nxb,
		otbnames: oxb.GetTbnames(),
		otables:  oxb.GetTables(),
	}
}

// 外部传入连接。用于正在使用数据库的情况，比如数据库运行期间备份数据库。
func NewDbback(oxb, nxb *Xb) *Dbback {
	return &Dbback{
		Oxb:      oxb,
		Nxb:      nxb,
		otbnames: oxb.GetTbnames(),
		otables:  oxb.GetTables(),
	}
}

// 按key,value导入，用于原、新数据库完全相同的情况。
func (d *Dbback) PutKV(k, v []byte) bool {
	err := d.Nxb.Db.Put(k, v, nil)
	if err != nil {
		//fmt.Println("err:", string(k), string(v))
		return false
	}
	//fmt.Println(string(k), string(v))
	return true
}

// 按字段数据导入，用于新数据库索引或字段名等改变的情况。
func (d *Dbback) PutData(k, v []byte) bool {
	if _, ok := d.otables[d.Otbname]; !ok {
		fmt.Println(d.Otbname, "表不存在。")
		return false
	}
	d.ntables = d.Nxb.GetTables()
	if _, ok := d.ntables[d.Ntbname]; !ok {
		fmt.Println(d.Ntbname, "表不存在。")
		return false
	}
	rdmap := d.otables[d.Otbname].RDtoMap(KVToRd(k, v, []int{}))
	for uk, uv := range d.Upfield { //修改字段名称，字段名称uk更名为uv
		rdmap[uv] = rdmap[uk]
		delete(rdmap, uk)
	}
	//return d.ntables[d.Ntbname].Ins(rdmap).Succ
	d.ntables[d.Ntbname].Ins(rdmap) //正向词
	d.PutDataR(k, v)
	return true
}

// 全文搜索反向排序，临时函数，未完善。
func (d *Dbback) PutDataR(k, v []byte) bool {
	//k=c-\x00\x00\x00\x01||\x00\x00\x00\x00
	bid := bytes.Replace(k, []byte(d.Ntbname), []byte(""), 1)
	rdmap := d.otables[d.Otbname].RDtoMap(KVToRd(k, v, []int{}))
	fidx := d.ntables[d.Ntbname].Ifo.FullText[0] //仅支持第一个字段的反向排序
	idx, _ := strconv.Atoi(fidx)
	fname := d.ntables[d.Ntbname].Ifo.Fields[idx]
	rs := rdmap[fname]
	text := []rune(rs)
	length := len(text)
	var result []rune
	for i := 0; i < length; i++ {
		result = append(result, text[length-i-1])
	}
	rdmap[fname] = string(result)
	kws := Analysis(rdmap[fname], []int{1})
	var nk []byte
	for _, v := range kws {
		if v == "" {
			continue
		}
		ftIdx := d.Nxb.GetTables()[d.Ntbname].ForDisparte(v, 7)
		for _, f := range ftIdx {
			nk = JoinBytes([]byte(d.Ntbname+IdxSplit+"r"+Split+f+Split), bid) ////c.r-如来-\x00\x00\x19a||\x00\x00\x00?
			d.Nxb.Db.Put(nk, []byte{}, nil)
		}
	}
	return true
}

// copy整个数据库
func (d *Dbback) CopyDb() {
	d.otables[d.otbnames[0]].Select.ForDbase(d.PutKV)
}

// copy一个表所有
func (d *Dbback) CopyTable(tbname string) {
	d.CopyTbItem(tbname)
	d.CopyTbInfo(tbname)
	d.CopyTbData(tbname)
}

// copy表tbname列表。即是将tbname加入表列表中。
func (d *Dbback) CopyTbItem(tbname string) {
	d.otables[tbname].Select.FindPrefixFun([]byte(Tbspfx+Split+tbname), true, d.PutKV)
}

// copy表信息
func (d *Dbback) CopyTbInfo(tbname string) {
	d.otables[tbname].Select.FindPrefixFun([]byte(TbInfopfx+Split+tbname+IdxSplit), true, d.PutKV)
}

// copy表数据
func (d *Dbback) CopyTbData(tbname string) {
	d.otables[tbname].Select.ForTb(d.PutKV)
}

// 仅仅copy一个表所有数据，索引等重新生成。
// 适应表明
func (d *Dbback) CopyTableData(tbname string) {
	d.CopyTbItem(tbname)
	d.CopyTbInfo(tbname)
	d.Otbname = tbname
	d.Ntbname = tbname
	d.ntables = d.Nxb.GetTables()
	d.Oxb.GetTables()[tbname].Select.ForRDFun(true, d.PutData)
}

// 获取新数据库表列表
func (d *Dbback) GetNewTbNames() {
	d.ntbnames = d.Nxb.GetTbnames()
}

// 获取新数据库所有表操作结构
func (d *Dbback) GetNewTables() {
	d.ntables = d.Nxb.GetTables()
}
