//小白数据库
//查找模板函数类
package xbdb

//游标根据开始位置进行顺序或倒序的遍历，并且合并kv为一个记录。
import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Iters struct {
	iter        iterator.Iterator //数据游标
	ok          bool              //游标是否结束
	asc         bool              //asc,升/降序
	bint, count int               //第几条开始
	//showFileds  []int             //要显示的字段对应的下标对应sql，select f0,f1,f2 from xxxx
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
func (i *Iters) ForData(showFileds []int, distinct bool) (r *TbData) {
	r = TbDatapool.Get().(*TbData)
	r.Release() //确保如果上次忘记释放内存不会导致数据混乱。
	loop, succ := -1, 0
	var rd []byte
	for i.ok {
		loop++
		if loop < i.bint {
			i.ok = itermove[i.asc](i.iter)
			continue
		}
		rd = KVToRd(i.iter.Key(), i.iter.Value(), showFileds)
		if distinct {
			if ArraysContains(r.Rd, rd) {
				i.ok = itermove[i.asc](i.iter)
				continue
			}
		}
		r.Rd = append(r.Rd, rd)
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
func ArraysContains(Rds [][]byte, Rd []byte) bool {
	for _, v := range Rds {
		if string(v) == string(Rd) {
			return true
		}
	}
	return false
}

//遍历KV结果集，执行函数为参数
//该函数主要用于遍历kv集。多个组合查询，都需要遍历k的结果集，并且进行交集处理。
func (i *Iters) ForKVFun(f func(k, v []byte) bool) {
	loop := 0
	for i.ok {
		loop++
		if f(i.iter.Key(), i.iter.Value()) {
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

//将key、value组合一个完整的记录，并提取需要的字段
//主要是用于select f1,f2,f3 from xxx获取对应的f1,f2,f3
func KVToRd(k, v []byte, idxs []int) (r []byte) {
	ks := SplitRd(k) //bytes.Split(k, []byte(Split))
	lks := len(ks)
	pk := ks[lks-1] //主键值必定是最后一个
	r = JoinBytes(SplitToCh(pk), []byte(Split), v)

	if len(idxs) == 0 { //select * from xxx
		return
	} else { //select f1,f5... from xxx
		ks = SplitRd(r)
		lks = len(ks)
		var idx int
		for i, v := range idxs {
			if v >= 0 {
				idx = v
			} else { //支持负数，即可取倒数第几个。主要用于取最后一个主键值的情况
				idx = lks + v
			}
			fv := SplitToCh(ks[idx])
			if i == 0 {
				r = fv
			} else {
				r = JoinBytes(r, []byte(Split), fv)
			}
		}
	}
	return
}

/*
//将key、value组合一个完整的记录，并提取需要的字段
func KVToRd(k, v []byte, idxs []int) (r []byte) {
	ks := SplitRd(k) //bytes.Split(k, []byte(Split))
	////将key去除第一个前缀，剩下的就是数据，有以下2中情况
	//索引：k=ca,fid-3-7 v=    得到:3-7-  (后面是空值)
	//主键: k=ca-1 v=ddd-ccdd-fff 得到:1-ddd-ccdd-fff

	ks = ks[1:]
	if len(idxs) == 0 { //select * from xxx
		for i, v := range ks { //转义
			ks[i] = SplitToCh(v)
		}
		r = JoinBytes(bytes.Join(ks, []byte(Split)), []byte(Split), v)
	} else { //select f1,f5... from xxx
		ks = append(ks, SplitRd(v)...)
		for i, v := range idxs {
			fv := SplitToCh(ks[v])
			if i == 0 {
				r = fv
			} else {
				r = JoinBytes(r, []byte(Split), fv)
			}
		}
	}
	return
}
*/

/*
//将key、value组合一个完整的记录
func KVToRdSf(k, v []byte) (r []byte) {
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

//获取索引key的主键id
func GetPKId(idxkey []byte) (r []byte) {
	r = GetIdxVal(idxkey, 0)

	return
}

//获取索引key的主键id
func GetIdxVal(rd []byte, idx int) (r []byte) {
	ks := SplitRd(rd)
	////将key去除第一个前缀，剩下的就是数据，有以下2中情况
	//索引：k=ca,fid-3-7 v=    得到:3-7-  (后面是空值)
	//主键: k=ca-1 v=ddd-ccdd-fff 得到:1-ddd-ccdd-fff
	r = ks[idx] //k的最后一个值
	return
}


//遍历记录结果集，执行函数为参数
//该函数主要用于遍历记录集
func (i *Iters) ForDataFun(f func(rd []byte) bool) {
	loop := 0
	for i.ok {
		loop++
		if f(i.KVToRd(i.iter.Key(), i.iter.Value())) {
			i.ok = itermove[i.asc](i.iter)
		} else {
			break
		}
	}
	i.iter.Release()
	//fmt.Println("跟踪执行次数:", loop) //跟踪执行次数
}

*/
