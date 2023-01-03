//小白数据库
//表信息
package xbdb

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

//表的类
type Table struct {
	Name   string
	Db     *leveldb.DB
	Select *Select
	Ac     *Autoinc
	Ifo    TableInfo
}

func NewTable(db *leveldb.DB, name string) *Table {
	return &Table{
		Name:   name,
		Db:     db,
		Select: NewSelect(name, db),
		Ifo:    NewTableInfo(db).GetInfo(name),
	}
}

//遍历分词
func (t *Table) ForDisparte(nr string, ftlen int) (disparte []string) {
	var knr string //, fid
	var ml, cl int
	var r, idxstr []rune
	r = []rune(nr)
	cl = len([]rune(nr))
	for cl > 0 {
		if cl >= ftlen {
			ml = ftlen
		} else {
			ml = cl
		}
		idxstr = r[:ml]
		knr = string(idxstr)
		disparte = append(disparte, knr)
		r = r[1:]
		cl = len(r)
	}
	return
}

var actmap map[string]func(k, v []byte) (r ReInfo)

//添加或删除一条记录，以及相关索引等所有数据
func (t *Table) Act(vals [][]byte, Act string) (r ReInfo) {
	if actmap == nil {
		actmap = map[string]func(k, v []byte) (r ReInfo){
			"insert": t.put,
			"delete": t.del,
		}
	}
	if len(vals) < len(t.Ifo.Fields) {
		r.Info = "字段参数长短不匹配！"
		return
	}
	/*
		//转义
		for i := 0; i < len(vals); i++ {
			vals[i] = t.Ifo.SplitToCh(vals[i])
		}*/
	r = t.ActPK(vals, Act)
	if !r.Succ {
		return
	}
	//添加表索引
	idx := -1
	var ivs []string
	idxfields := ""
	var idxvals []byte
	for _, iv := range t.Ifo.Idxs {
		if iv == "" {
			continue
		}
		ivs = strings.Split(iv, ",")
		for i := 0; i < len(ivs); i++ { //组织单个或组合索引key
			idx, _ = strconv.Atoi(ivs[i])
			idxfields += t.Ifo.Fields[idx]
			idxvals = JoinBytes(idxvals, vals[idx])
			if i != len(ivs)-1 {
				idxfields += IdxSplit
				idxvals = JoinBytes(idxvals, []byte(IdxSplit))
			}
		}
		r = t.ActIDX([]byte(idxfields), idxvals, vals[0], []byte{}, Act)
		if !r.Succ {
			return
		}
		idxfields = ""
		idxvals = idxvals[:0]
	}
	//添加表全文索引
	ftlen, _ := strconv.Atoi(t.Ifo.FTLen)
	var ftIdx []string
	for _, i := range t.Ifo.FullText {
		if i == "" {
			continue
		}
		idx, _ = strconv.Atoi(i)
		ftIdx = t.ForDisparte(string(vals[idx]), ftlen)
		for p, f := range ftIdx {
			t.ActIDX([]byte(t.Ifo.Fields[idx]), []byte(f), vals[0], IntToBytes(p), Act)
			if !r.Succ {
				return
			}
		}
	}
	r.Succ = true
	r.Info = "成功！"
	return
}

//添加/删除主键数据，即添加/删除一条记录。
func (t *Table) ActPK(vals [][]byte, Act string) (r ReInfo) {
	key := t.Select.GetPkKey(vals[0])
	r = actmap[Act](key, bytes.Join(vals[1:], []byte(Split)))
	if !r.Succ {
		return
	}
	return
}

//添加/删除一条索引数据。
func (t *Table) ActIDX(idxfield, idxvalue, pkvalue, val []byte, Act string) (r ReInfo) {
	//bySplit := []byte(Split)
	//k=ca,fid-3-7 v=
	//prefix := JoinBytes([]byte(t.Ifo.Name+","), idxFieldname, bySplit, idxFieldvalue, bySplit, PKvalue)
	key := t.Select.GetIdxPrefixKey(idxfield, idxvalue, pkvalue) //getIdxPrefixKey
	r = actmap[Act](key, val)
	if !r.Succ {
		return
	}
	return
}

//将记录转换为map
func (t *Table) RDtoMap(Rd []byte, ifo *TableInfo) (r map[string]string) {
	r = make(map[string]string, len(ifo.Fields))
	vs := bytes.Split(Rd, []byte(Split))
	for i, v := range vs {
		r[ifo.Fields[i]] = ifo.ByteChString(ifo.FieldType[i], v)
	}
	return
}

/*
前缀遍历
bint 第几条开始
asc,升/降序

func (t *Table) FindPrefix(key []byte, asc bool, b, count int) (r *TbDataMap) {
	iter, ok := t.Select.IterPrefixMove(key, asc)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForDataToMap(&t.Ifo)
	return
}


范围遍历
asc,升/降序

func (t *Table) FindRand(bkey, ekey []byte, asc bool, b, count int) (r *TbDataMap) {
	iter, ok := t.Select.IterRandMove(bkey, ekey, asc)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForDataToMap(&t.Ifo)
	return
}


定位遍历
asc,升/降序

func (t *Table) FindSeek(key []byte, asc bool, b, count int) (r *TbDataMap) {
	iter, ok := t.Select.IterSeekMove(key)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForDataToMap(&t.Ifo)
	return
}
*/
