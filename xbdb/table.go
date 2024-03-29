// 小白数据库
// 表信息
package xbdb

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var mu sync.RWMutex
var actmap map[string]func(k, v []byte) (r ReInfo)

// 表的类
type Table struct {
	Name   string
	db     *leveldb.DB
	Select *Select
	Ac     *Autoinc
	Ifo    *TableInfo
}

func NewTable(db *leveldb.DB, name string) *Table {
	return &Table{
		Name:   name,
		db:     db,
		Select: NewSelect(db, name),
		Ifo:    NewTableInfo(db, name),
	}
}

// 遍历分词
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

// 添加或删除一条记录，以及相关索引等所有数据
func (t *Table) Act(vals [][]byte, Act string) (r ReInfo) {
	return t.Acts(vals, Act, nil)
}

// 添加或删除一条记录，以及相关索引等所有数据等事务
// updatefield,修改时用。用于记录需要修改那个字段。与表字段一一对应。
// 修改某个某些字段时，updatefield作为判断，不用把所有索引都删除再重新添加，导致性能不高和不灵活。
func (t *Table) Acts(vals [][]byte, Act string, updatefield []bool) (r ReInfo) {
	if actmap == nil {
		actmap = map[string]func(k, v []byte) (r ReInfo){
			"insert": t.put,
			"delete": t.del,
		}
	}
	r = t.ActPK(vals, Act)
	if !r.Succ {
		return
	}
	r = t.AddIdx(vals, Act, updatefield) //添加普通索引
	if !r.Succ {
		return
	}
	r = t.AddFullIdx(vals, Act, updatefield) //添加全文索引
	return
}
func (t *Table) AddIdx(vals [][]byte, Act string, updatefield []bool) (r ReInfo) {
	//添加表索引
	idx := -1
	var ivs []string
	idxfields := ""
	var idxval, idxvals []byte
	for _, iv := range t.Ifo.Idxs {
		if iv == "" {
			//r.Succ = true
			//r.Info = "没有索引"
			continue
		}
		ivs = strings.Split(iv, ",")
		if updatefield != nil { //修改的情况
			if !isUpdateField(ivs, updatefield) {
				continue //不是修改字段则退出，不用添加或删除原有索引
			}
		}
		for i := 0; i < len(ivs); i++ { //组织单个或组合索引key
			idx, _ = strconv.Atoi(ivs[i])
			idxval = vals[idx]
			if len(idxval) == 0 { //不添加空值或nil值的索引。利大于弊。
				r.Succ = true
				return
			}
			idxfields += t.Ifo.Fields[idx]       //累加
			idxvals = JoinBytes(idxvals, idxval) //累加
			if i != len(ivs)-1 {                 //不是末尾，则加分隔符
				idxfields += IdxSplit
				idxvals = JoinBytes(idxvals, []byte(IdxSplit))
			}
		}
		r = t.ActIDX([]byte(idxfields), idxvals, vals[0], []byte{}, Act)
		if !r.Succ {
			return
		}
		//重置
		idxfields = ""
		idxvals = idxvals[:0]
	}
	r.Succ = true
	r.Info = "成功。"
	return
}
func (t *Table) AddFullIdx(vals [][]byte, Act string, updatefield []bool) (r ReInfo) {
	idx := -1
	ftlen, _ := strconv.Atoi(t.Ifo.FTLen)
	var ftIdx []string
	for _, i := range t.Ifo.FullText {
		if i == "" {
			continue
		}
		idx, _ = strconv.Atoi(i)
		if updatefield != nil { //修改的情况
			if !updatefield[idx] { //非修改字段不添加/删除索引
				continue
			}
		}
		var patterns []int //搜索词解析模型
		var iv int
		for _, v := range t.Ifo.Patterns { //字符串数组转int数组
			iv, _ = strconv.Atoi(v)
			patterns = append(patterns, iv)
		}
		kws := Analysis(string(vals[idx]), patterns)
		for i, v := range kws {
			if v == "" {
				continue
			}
			ftIdx = t.ForDisparte(v, ftlen)
			for p, f := range ftIdx {
				r = t.ActIDX([]byte(t.Ifo.Fields[idx]), []byte(f), vals[0], IntToBytes(p+i), Act)
				if !r.Succ {
					return
				}
			}
		}
	}
	r.Succ = true
	r.Info = "成功。"
	return
}

// 是不是修改字段。
// 由于支持组合索引，故而需要循环，看似复制些
func isUpdateField(ks []string, updatefield []bool) bool {
	isf := false
	idx := 0
	for i := 0; i < len(ks); i++ { //单个或组合索引。由于支持组合索引，故而需要循环，看似复制些
		idx, _ = strconv.Atoi(ks[i])
		isf = isf || updatefield[idx]
	}
	return isf
}

// 添加/删除主键数据，即添加/删除一条记录。
func (t *Table) ActPK(vals [][]byte, Act string) (r ReInfo) {
	key := t.Select.GetPkKey(vals[0])
	r = actmap[Act](key, bytes.Join(vals[1:], []byte(Split)))
	if !r.Succ {
		return
	}
	return
}

// 添加/删除一条索引数据。
func (t *Table) ActIDX(idxfield, idxvalue, pkvalue, val []byte, Act string) (r ReInfo) {
	/*
		if len(idxvalue) == 0 { //len(idxvalue) == 0兼容==nil的情况。
			r.Succ = true
			r.Info = "不对nil或空值创建索引。"
			return //不对nil或空值创建索引。利大于弊。
		}*/
	key := t.Select.GetIdxPrefixKey(idxfield, idxvalue, pkvalue)
	r = actmap[Act](key, val)
	/*
		if !r.Succ {
			return
		}
	*/
	return
}

// 字符串转byte
// params字段对应的字符串map
func (t *Table) StrToByte(params map[string]string) (r [][]byte) {
	for i, v := range t.Ifo.Fields {
		r = append(r, t.Ifo.TypeChByte(t.Ifo.FieldType[i], params[v]))
	}
	return
}

/*
//字符串根据表类型信息转换为byte数据
func StrToByteForIfo(params map[string]string, ifo TableInfo) (r [][]byte) {
	for i, v := range ifo.Fields {
		r = append(r, ifo.TypeChByte(ifo.FieldType[i], params[v]))
	}
	return
}
*/
//将记录转换为map
func (t *Table) RDtoMap(Rd []byte) (r map[string]string) {
	r = t.FieldValuetoMap(Rd, t.Ifo)
	return
}

// 将记录转换为map
func (t *Table) FieldValuetoMap(Rd []byte, Ifo *TableInfo) (r map[string]string) {
	r = make(map[string]string, len(Ifo.Fields))
	vs := t.Split(Rd)
	vslen := len(vs)
	for i, v := range Ifo.Fields {
		if i < vslen {
			r[v] = Ifo.ByteChString(Ifo.FieldType[i], vs[i])
		} else { //添加新字段时，字段比旧数据数组长度大。
			r[v] = ""
		}
	}
	/*
		for i, v := range vs {
			r[Ifo.Fields[i]] = Ifo.ByteChString(Ifo.FieldType[i], v) //将包括分隔符的转义数据恢复
		}*/
	return
}

// 将记录分开并转义数据恢复
func (t *Table) Split(Rd []byte) (r [][]byte) {
	r = SplitRd(Rd)
	return
}

var Bufpool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

/*
	func (t *Table) DataToJson(tbd *TbData) (r *bytes.Buffer) {
		r = t.DataToJsonforIfo(tbd, t.Ifo)
		return
	}
*/
func (t *Table) DataToJsonApp(tbd *TbData) (r *bytes.Buffer) {
	r = t.DataToJsonforIfo(tbd, t.Ifo)
	return
}

func (t *Table) DataToJsonforIfo(tbd *TbData, Ifo *TableInfo) (r *bytes.Buffer) {
	if tbd == nil {
		return
	}
	r = Bufpool.Get().(*bytes.Buffer)
	if r.Len() > 0 { //保证数据不混乱
		r.Reset()
	}
	var rdmap map[string]string
	jsonstr := ""
	valstr := ""
	r.WriteString("{\"result\":[")
	for j, v := range tbd.Rd {
		if v == nil {
			continue
		}
		r.WriteString("{")
		rdmap = t.FieldValuetoMap(v, Ifo)
		/*
			[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
			{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
		*/
		for i, fv := range Ifo.FieldType {
			switch fv {
			case "string", "time", "bool":
				valstr = strconv.Quote(rdmap[Ifo.Fields[i]])
			default:
				valstr = rdmap[Ifo.Fields[i]]
			}
			jsonstr = "\"" + Ifo.Fields[i] + "\":" + valstr
			if i != len(rdmap)-1 {
				jsonstr += ","
			}
			r.WriteString(jsonstr)
		}
		r.WriteString("}")
		if j != len(tbd.Rd)-1 {
			r.WriteString(",")
		}
	}
	r.WriteString("]}")
	tbd.Release()
	return
}

/*
func (t *Table) DataToJsonforIfoApp(tbd *TbData, Ifo *TableInfo) (r *bytes.Buffer) {
	if tbd == nil {
		return
	}
	r = Bufpool.Get().(*bytes.Buffer)
	if r.Len() > 0 { //保证数据不混乱
		r.Reset()
	}
	var rdmap map[string]string
	jsonstr := ""
	valstr := ""
	r.WriteString("[")
	for j, v := range tbd.Rd {
		if v == nil {
			continue
		}
		r.WriteString("{")
		rdmap = t.FieldValuetoMap(v, Ifo)

		for i, fv := range Ifo.FieldType {
			switch fv {
			case "string", "time", "bool":
				valstr = strconv.Quote(rdmap[Ifo.Fields[i]])
			default:
				valstr = rdmap[Ifo.Fields[i]]
			}
			jsonstr = "\"" + Ifo.Fields[i] + "\":" + valstr
			if i != len(rdmap)-1 {
				jsonstr += ","
			}
			r.WriteString(jsonstr)
		}
		r.WriteString("}")
		if j != len(tbd.Rd)-1 {
			r.WriteString(",")
		}
	}
	r.WriteString("]")
	tbd.Release()
	return
}
*/
// 获取字段在表中的索引id
func (t *Table) GetFieldIdx(field string) int {
	for i, fv := range t.Ifo.Fields {
		if field == fv { //得到pv在Fields中索引id
			return i
		}
	}
	return -1
}

// 根据字段索引判断是否索引
func (t *Table) FieldIsIdx(idx int) bool {
	var i int
	if idx == -1 {
		return false
	}
	for _, v := range t.Ifo.Idxs {
		i, _ = strconv.Atoi(v)
		if idx == i {
			return true
		}
	}
	return false
}
