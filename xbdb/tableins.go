//小白数据库
//表信息
package xbdb

import "bytes"

//添加一条数据以及相关索引等所有数据，默认数据与字段一一对应。
func (t *Table) Insert(vals [][]byte) (r ReInfo) {
	if string(vals[0]) == "\x00\x00\x00\x00" { //如果主键为空，则是使用自动增值
		if t.Ac == nil {
			t.newAutoinc()
		}
		vals[0] = t.Ifo.TypeChByte("int", t.Ac.GetidStr())
	}
	r = t.Act(vals, "insert")
	return
}

//添加一条数据以及相关索引等所有数据，默认数据与字段一一对应。
func (t *Table) Ins(params map[string]string) (r ReInfo) {
	if string(params["id"]) == "" { //如果主键为空，则是使用自动增值
		if t.Ac == nil {
			t.newAutoinc()
		}
		params["id"] = t.Ac.GetidStr()
	}
	vals := t.StrToByte(params)
	r = t.Act(vals, "insert")
	return
}

//添加一个kv
func (t *Table) put(k, v []byte) (r ReInfo) {
	err = t.Db.Put(k, v, nil) //vals[0]=主键
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
	iter := t.Select.IterPrefix([]byte(t.Name + Split))
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

/*
//添加主键数据，即添加一条记录。
func (t *Table) InsPK(vals [][]byte) (r ReInfo) {
	//prefix := t.Ifo.Name + Split
	//ca-7 禅定品-3-1
	//k=ca-7 v=禅定品-3-1
	key := t.Ifo.getPkKey(vals[0])
	r = t.put(key, bytes.Join(vals[1:], []byte(Split)))
	if !r.Succ {
		return
	}
	return
}

//添加一条索引数据。
func (t *Table) InsIDX(idxfield, idxvalue, pkvalue []byte) (r ReInfo) {
	//bySplit := []byte(Split)
	//k=ca,fid-3-7 v=
	//prefix := JoinBytes([]byte(t.Ifo.Name+","), idxFieldname, bySplit, idxFieldvalue, bySplit, PKvalue)
	key := t.Ifo.getIdxPfxKey(idxfield, idxvalue, pkvalue)
	r = t.put(key, []byte{}) //vals[0]=主键
	if !r.Succ {
		return
	}
	return
}
*/
