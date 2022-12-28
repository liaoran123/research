//小白数据库
//表信息
package xbdb

import (
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	Split      = "--"    //字段分隔符
	ChSplit    = "#f0"   //字段分隔符的转义码
	IdxSplit   = ".."    //索引分隔符，转义码 #f1
	ChIdxSplit = "#f1"   //索引分隔符的转义码
	TbInfopfx  = "tbifo" //表信息的前缀
	Tbspfx     = "table" //表列表的前缀
)

//表信息的类。//默认必须第一个字段是主键id
type TableInfo struct {
	Db        *leveldb.DB
	Name      string   //表名
	Fields    []string //字段
	FieldType []string //字段对应的类型
	Pk        string   //默认必须有一个自动增值的主键id
	Idxs      []string //索引字段的下标，不使用[]int，转换byte太麻烦。
	FullText  []string //考据级全文搜索索引字段的下标。
	FTLen     string   //全文搜索的长度，中文默认是7
}

func NewTableInfo(DB *leveldb.DB) *TableInfo {
	return &TableInfo{
		Db: DB,
	}
}

//创建/修改一个表，默认第一个字段必须是主键
func (t *TableInfo) Create(name, ftlen string, fields, fieldType, idxs, fullText []string) (r ReInfo) {
	if len(fieldType) != len(fields) {
		r.Info = "字段和类型数据不匹配！"
		return
	}

	tbpfx := TbInfopfx + Split + name //表信息前缀

	r.Succ = t.Db.Put([]byte(tbpfx), []byte(strings.Join(fields, Split)), nil) == nil                            //添加字段信息
	r.Succ = r.Succ && t.Db.Put([]byte(tbpfx+IdxSplit+"ty"), []byte(strings.Join(fieldType, Split)), nil) == nil //添加字段类型信息
	r.Succ = r.Succ && t.Db.Put([]byte(tbpfx+IdxSplit+"pk"), []byte(fields[0]), nil) == nil                      //添加主键信息
	r.Succ = r.Succ && t.Db.Put([]byte(tbpfx+IdxSplit+"idx"), []byte(strings.Join(idxs, Split)), nil) == nil     //添加索引信息
	r.Succ = r.Succ && t.Db.Put([]byte(tbpfx+IdxSplit+"ft"), []byte(strings.Join(fullText, Split)), nil) == nil  //添加索引信息
	r.Succ = r.Succ && t.Db.Put([]byte(tbpfx+IdxSplit+"ftlen"), []byte(ftlen), nil) == nil

	r.Succ = r.Succ && t.Db.Put([]byte(Tbspfx+Split+name), []byte{}, nil) == nil //添加表列表
	if r.Succ {
		r.Info = "创建表“" + name + "”成功！"
	} else {
		r.Info = "创建表“" + name + "”失败！"
	}
	return
}

func (t *TableInfo) GetInfo(name string) (tbif TableInfo) {

	tbpfx := TbInfopfx + Split + name //表信息前缀
	tf := TableInfo{}
	tf.Name = t.Name
	var data []byte
	var err error
	data, err = t.Db.Get([]byte(tbpfx), nil) //打开字段信息
	if err != nil {
		return
	}
	tf.Fields = strings.Split(string(data), Split)
	data, _ = t.Db.Get([]byte(tbpfx+IdxSplit+"ty"), nil) //打开主键信息
	tf.FieldType = strings.Split(string(data), Split)
	data, _ = t.Db.Get([]byte(tbpfx+IdxSplit+"pk"), nil) //打开主键信息
	tf.Pk = string(data)
	data, _ = t.Db.Get([]byte(tbpfx+IdxSplit+"idx"), nil) //打开索引信息
	tf.Idxs = strings.Split(string(data), Split)
	data, _ = t.Db.Get([]byte(tbpfx+IdxSplit+"ft"), nil) //打开全文索引信息
	tf.FullText = strings.Split(string(data), Split)
	data, _ = t.Db.Get([]byte(tbpfx+IdxSplit+"ftlen"), nil) //全文索引长度信息
	tf.FTLen = string(data)

	tbif = tf
	return
}

//删除表信息
func (t *TableInfo) Del(name string) (r ReInfo) {
	if name == "" {
		return
	}
	tbpfx := TbInfopfx + Split + name //表信息前缀
	iter := t.Db.NewIterator(util.BytesPrefix([]byte(tbpfx)), nil)
	for iter.Next() {
		err = t.Db.Delete(iter.Key(), nil)
		if err != nil {
			r.Info = err.Error()
			r.Succ = false
			break
		}
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		r.Info = err.Error()
		r.Succ = false
	}
	r.Info = "删除" + name + "成功！"
	r.Succ = true
	return
}

//将字段类型数据转换为对应的[]byte数据数组
func (t *TableInfo) TypeChByte(FieldType, Fieldvalue string) (r []byte) {
	switch FieldType {
	case "int":
		iv, _ := strconv.Atoi(Fieldvalue)
		r = IntToBytes(iv)
	case "int64":
		iv, _ := strconv.Atoi(Fieldvalue)
		r = Int64ToBytes(int64(iv))
	case "float":
		fv, _ := strconv.ParseFloat(Fieldvalue, 64) //只能转Float64
		r = Float64ToByte(fv)

	default:
		r = []byte(Fieldvalue)
	}
	return
}

//根据获取字段对应的类型
func (t *TableInfo) GetFieldType(idxfield string) (r string) {
	for i, v := range t.Fields {
		if v == idxfield {
			return t.FieldType[i]
		}
	}
	return
}

//将字段名称数据转换为对应的[]byte数据数组
func (t *TableInfo) FieldChByte(field, value string) (r []byte) {
	filedtype := t.GetFieldType(field)
	return t.TypeChByte(filedtype, value)
}

//将字段[]byte数据数组转换为字符串
func (t *TableInfo) ByteChString(FieldType string, val []byte) (r string) {
	switch FieldType {
	case "int":
		ir := BytesToInt(val)
		r = strconv.Itoa(ir)
	case "int64":
		ir := BytesToInt64(val)
		r = strconv.FormatInt(ir, 10)
	case "float":
		ir := ByteToFloat64(val)
		r = strconv.FormatFloat(ir, 'f', 2, 64) //2,2位小数点
	default:
		r = string(val)
	}
	return
}

//将表字段[]byte数据数组转换为字符串数组
func (t *TableInfo) ValsChString(vals [][]byte) (r []string) {
	cs := ""
	for i, v := range vals {
		cs = t.ByteChString(t.FieldType[i], v)
		r = append(r, cs)
	}
	return
}

/*
//将包括分隔符的数据转义
func (t *TableInfo) SplitToCh(k []byte) (r []byte) {
	r = bytes.Replace(k, []byte(Split), []byte(ChSplit), -1)
	r = bytes.Replace(r, []byte(IdxSplit), []byte(ChIdxSplit), -1)
	return
}

//将包括分隔符的转义数据恢复
func (t *TableInfo) ChToSplit(k []byte) (r []byte) {
	r = bytes.Replace(k, []byte(ChSplit), []byte(Split), -1)
	r = bytes.Replace(r, []byte(ChIdxSplit), []byte(IdxSplit), -1)
	return
}
*/
