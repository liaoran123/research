// 小白数据库
// 表信息
package xbdb

import (
	"bytes"
	"fmt"
)

// 删除一条数据记录以及相关索引等所有数据
// pk=id值
func (t *Table) Delete(pk []byte) (r ReInfo) {
	//key := JoinBytes([]byte(t.Name+Split), pk)
	key := JoinBytes(t.Select.GetTbKey(), pk)
	data, err := t.db.Get(key, nil)
	if err != nil {
		fmt.Println(key)
		r.Info = err.Error()
		return
	}
	vals := bytes.Split(data, []byte(Split))
	//将pk在vals前面插入，以便对应索引的下标
	vals = append(vals, []byte{}) // 切片扩展1个空间
	copy(vals[0+1:], vals[0:])    // a[i:]向后移动1个位置
	vals[0] = pk

	//删除表记录
	r = t.Act(vals, "delete")
	return
}

// 入口参数为字符串的删除函数
func (t *Table) Del(pk string) (r ReInfo) {
	bpk := t.Ifo.TypeChByte(t.Ifo.FieldType[0], pk)
	r = t.Delete(bpk)
	return
}

// 按删除一个key
func (t *Table) del(k, v []byte) (r ReInfo) {
	mu.Lock()
	err = t.db.Delete(k, nil)
	mu.Unlock()
	if err != nil {
		r.Info = err.Error()
		r.Succ = false
	} else {
		r.Info = "删除成功！"
		r.Succ = true
	}
	return
}

// 按删除一个key，测试使用
func (t *Table) DelTest(k, v []byte) (r ReInfo) {
	mu.Lock()
	err = t.db.Delete(k, nil)
	mu.Unlock()
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

//删除子索引
key=
ca,fid-2-1
ca,fid-2-5
ca,fid-2-6

*/
