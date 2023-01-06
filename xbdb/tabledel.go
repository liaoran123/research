//小白数据库
//表信息
package xbdb

import (
	"bytes"
)

//删除一条数据以及相关索引等所有数据
//pk=id值
func (t *Table) Delete(pk []byte) (r ReInfo) {
	key := JoinBytes([]byte(t.Name+Split), pk)
	data, err := t.Db.Get(key, nil)
	if err != nil {
		r.Info = err.Error()
		return
	}
	vals := bytes.Split(data, []byte(Split))
	//将pk在vals前面插入，以便对应索引的下标
	vals = append(vals, []byte{}) // 切片扩展1个空间
	copy(vals[0+1:], vals[0:])    // a[i:]向后移动1个位置
	vals[0] = pk
	/*               // 设置新添加的元素
	//反转义数据
	for i := 0; i < len(vals); i++ {
		vals[i] = t.Ifo.ChToSplit(vals[i])
	}*/
	//删除表索引
	r = t.Act(vals, "delete")
	return
}

//入口参数为字符串的删除函数
func (t *Table) Del(pk string) (r ReInfo) {
	bpk := t.Ifo.TypeChByte(t.Ifo.FieldType[0], pk)
	r = t.Delete(bpk)
	return
}

//按删除一个key
func (t *Table) del(k, v []byte) (r ReInfo) {
	err = t.Db.Delete(k, nil)
	if err != nil {
		r.Info = err.Error()
		r.Succ = false
	} else {
		r.Info = "删除成功！"
		r.Succ = true
	}
	return
}

/*
ca-1 乾隆大藏经-0-0
ca-2 金刚经-1-0
ca-3 六祖坛经-1-0
ca-4 机缘品-3-1
ca-5 般若品-3-1
//根据主键pk删除表的一条记录数据。key=ca-2

func (t *Table) DelPK(pk []byte) (r ReInfo) {
	//tbpfx := t.Ifo.Name + Split
	//bpfx := JoinBytes([]byte(tbpfx), pk) //key=ca-2
	bpfx := t.Ifo.getPkKey(pk)
	r = t.delpfx(bpfx)
	return
}
*/
/*
//删除子索引
key=
ca,fid-2-1
ca,fid-2-5
ca,fid-2-6

//根据删除一个索引的所有数据.key=ca,fid-2
func (t *Table) DelIDX(idxfield, idxvalue []byte) (r ReInfo) {
	bpfx := t.Select.getIdxPfx(idxfield, idxvalue) //t.Ifo.getIdxPfx(idxfield, idxvalue)
	r = t.delpfx([]byte(bpfx))
	return
}


key=
ca,fid-2-1
ca,fid-2-5
ca,fid-2-6

根据删除一个索引的一个数据.key=ca,fid-2-1

func (t *Table) DelIDXPK(idxfield, idxvalue, pkvalues []byte) (r ReInfo) {
	//bSplit := []byte(Split)
	//bpfx := JoinBytes([]byte(t.Ifo.Name+","), idxfield, bSplit, idxvalue, bSplit, pkvalues)
	bpfx := t.Ifo.getIdxPfxKey(idxfield, idxvalue, pkvalues)
	r = t.del(bpfx, []byte{})
	return
}

//删除整个表
func (t *Table) DelAll() (r ReInfo) {
	bpfx := t.Name
	r = t.delpfx([]byte(bpfx))
	return
}

//按前缀删除数据
func (t *Table) delpfx(tbpfx []byte) (r ReInfo) {
	iter := t.Db.NewIterator(util.BytesPrefix(tbpfx), nil)
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
	} else {
		r.Info = "删除成功！"
		r.Succ = true
	}
	return
}
*/
