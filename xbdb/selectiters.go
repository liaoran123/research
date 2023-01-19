//小白数据库
//查找模板函数类
package xbdb

import (
	"bytes"
	"sync"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Iters struct {
	iter        iterator.Iterator //数据游标
	ok          bool              //游标是否结束
	asc         bool              //asc,升/降序
	bint, count int               //第几条开始
}
type TbData struct {
	Rd [][]byte //记录
}

func (t *TbData) Release() {
	if t == nil {
		return
	}
	t.Rd = t.Rd[:0]
	TbDatapool.Put(t)
}

var TbDatapool = sync.Pool{
	New: func() interface{} {
		return new(TbData)
	},
}

func NewIters(iter iterator.Iterator, ok, asc bool, bint, count int) *Iters {
	return &Iters{
		iter:  iter,
		ok:    ok,
		asc:   asc,
		bint:  bint,
		count: count,
	}
}

//遍历结果集模板函数
func (i *Iters) ForData() (r *TbData) {
	r = TbDatapool.Get().(*TbData)
	loop, succ := -1, 0
	for i.ok {
		loop++
		if loop < i.bint {
			i.ok = itermove[i.asc](i.iter)
			continue
		}
		r.Rd = append(r.Rd, KVToRd(i.iter.Key(), i.iter.Value()))
		succ++
		i.ok = itermove[i.asc](i.iter)
		if i.count != -1 { //-1不限制条数
			if succ >= i.count {
				break
			}
		}
	}
	i.iter.Release()
	return
}

//遍历结果集，执行函数为参数
func (i *Iters) ForDataFun(f func(rd []byte) bool) {
	loop := 0
	for i.ok {
		loop++
		if f(KVToRd(i.iter.Key(), i.iter.Value())) {
			i.ok = itermove[i.asc](i.iter)
		} else {
			break
		}
	}
	i.iter.Release()
	//fmt.Println("跟踪执行次数:", loop) //跟踪执行次数
}

//遍历累计记录数
func (i *Iters) ForDataCount() (r int) {
	for i.ok {
		r++
		i.ok = itermove[i.asc](i.iter)
	}
	i.iter.Release()
	return
}

//将key、value组合一个完整的记录
func KVToRd(k, v []byte) (r []byte) {
	ks := SplitRd(k) //bytes.Split(k, []byte(Split))
	////将key去除第一个前缀，剩下的就是数据，有以下2中情况
	//索引：k=ca,fid-3-7 v=    得到:3-7-  (后面是空值)
	//主键: k=ca-1 v=ddd-ccdd-fff 得到:1-ddd-ccdd-fff
	ks = ks[1:]
	for i, v := range ks { //转义
		ks[i] = SplitToCh(v)
	}
	r = JoinBytes(bytes.Join(ks, []byte(Split)), []byte(Split), v)
	return
}

/*
var TbDataMappool = sync.Pool{
	New: func() interface{} {
		return new(TbDataMap)
	},
}

type TbDataMap struct {
	Rd []map[string]string //记录
}

func (t *TbDataMap) Release() {
	if t == nil {
		return
	}
	t.Rd = t.Rd[:0]
	TbDataMappool.Put(t)
}
//遍历结果集，转换为[]map[string]string
func (i *Iters) ForDataToMap(ifo *TableInfo) (r *TbDataMap) { //[]map[string]string
	r = TbDataMappool.Get().(*TbDataMap)
	var ks [][]byte
	var key []byte
	for i.ok {
		tmap := make(map[string]string, len(ifo.Fields)) //必须放在循环内初始化，否则值不会改变
		ks = bytes.Split(i.iter.Key(), []byte(Split))
		key = ks[len(ks)-1]
		tmap[ifo.Fields[0]] = ifo.ByteChString(ifo.FieldType[0], key)
		ks = bytes.Split(i.iter.Value(), []byte(Split))
		for j, v := range ks {
			tmap[ifo.Fields[j+1]] = ifo.ByteChString(ifo.FieldType[j+1], v)
		}
		r.Rd = append(r.Rd, tmap)
		i.ok = itermove[i.asc](i.iter)
	}
	i.iter.Release()
	return
}

*/
