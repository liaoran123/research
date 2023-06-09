// 小白数据库
// 表信息
package xbdb

import (
	"bytes"
	"strconv"
)

// 添加一条数据以及相关索引等所有数据，默认数据与字段一一对应。
func (t *Table) Insert(vals [][]byte) (r ReInfo) {
	//if string(vals[0]) == "" { //"\x00\x00\x00\x00" { //如果主键为空，则是使用自动增值。只有类型是int时成立。
	if t.Ac == nil {
		t.newAutoinc()
	}
	if len(vals[0]) == 0 {
		vals[0] = t.Ifo.TypeChByte("int", t.Ac.GetidStr())
		r.LastId = strconv.Itoa(t.Ac.GetidDic())
	} else {
		t.Ac.id = BytesToInt(vals[0]) + 1 //将用户提交的id+1设置为自动增值的最后id。
	}
	r = t.Act(vals, "insert")
	if r.Succ {
		if r.LastId == "" { //非自动增值的情况
			r.LastId = string(vals[0])
		}
	}
	return
}

// 添加一条数据以及相关索引等所有数据，默认数据与字段一一对应。
func (t *Table) Ins(params map[string]string) (r ReInfo) {
	if t.Ac == nil {
		t.newAutoinc()
	}
	if _, ok := params["id"]; !ok { //如果主键为空，则是使用自动增值,只有类型是int时成立。

		params["id"] = t.Ac.GetidStr()
	} else {
		t.Ac.id, _ = strconv.Atoi(params["id"])
		t.Ac.id++ //将用户提交的id+1设置为自动增值的最后id
	}
	vals := t.StrToByte(params)
	r = t.Act(vals, "insert")
	if r.Succ {
		r.LastId = params["id"]
	}
	return
}

// 添加一个kv
func (t *Table) put(k, v []byte) (r ReInfo) {
	mu.Lock()
	err = t.db.Put(k, v, nil) //vals[0]=主键
	mu.Unlock()
	if err != nil {
		r.Succ = false
		r.Info = err.Error()
		return
	}
	r.Succ = true
	r.Info = "put成功！"
	return
}
func (t *Table) newAutoinc() {
	iter := t.Select.IterPrefix(t.Select.GetTbKey()) //[]byte(t.Name + Split)
	var key []byte
	if iter.Last() {
		key = iter.Key()
	}
	ks := bytes.Split(key, []byte(Split))
	bid := ks[len(ks)-1]

	id := 1
	if bid != nil {
		id = BytesToInt(bid) + 1
	}
	t.Ac = NewAutoinc(id)
}
