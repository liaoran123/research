// 小白数据库
// 表信息
package xbdb

//根据getprefix.go和getiters.go，也即key和游标的各种组合。

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

/*
查询执行流程，
1，根据getprefix.go的表前缀规则获取key
2，通过key配合getiters.go各个函数获取各种查询游标数据iters初始化（主要是按索引和主键查询和顺序和倒序）。
3，根据itersfor.go进行各种遍历。
*/
type Select struct {
	db     *leveldb.DB
	Tbname string
}

var (
	iterFixed map[bool]func(iter iterator.Iterator) bool //起始位置。升序，first；降序，last
	itermove  map[bool]func(iter iterator.Iterator) bool //移动netx, prev
)

// 下面四个函数为了动态的顺序和倒序的遍历游标
func First(iter iterator.Iterator) bool {
	return iter.First()
}
func Last(iter iterator.Iterator) bool {
	return iter.Last()
}
func Prev(iter iterator.Iterator) bool {
	return iter.Prev()
}
func Next(iter iterator.Iterator) bool {
	return iter.Next()
}

func NewSelect(db *leveldb.DB, tbname string) *Select { //*leveldb.DB {

	iterFixed = make(map[bool]func(iter iterator.Iterator) bool, 2)
	iterFixed[true] = First
	iterFixed[false] = Last
	itermove = make(map[bool]func(iter iterator.Iterator) bool, 2)
	itermove[true] = Next
	itermove[false] = Prev

	return &Select{
		db:     db,
		Tbname: tbname,
	}
}

// 遍历数据库，主要用于复制数据库
func (s *Select) ForDbase(f func(k, v []byte) bool) {
	iter := s.Nil()
	for iter.Next() {
		if f(iter.Key(), iter.Value()) {
		} else {
			iter.Release()
			return
		}
	}
	iter.Release()
}

// 遍历表所有，主要用于复制
func (s *Select) ForTb(f func(k, v []byte) bool) {
	iter := s.IterPrefix(s.GetTbLikeKey())
	for iter.Next() {
		if f(iter.Key(), iter.Value()) {
		} else {
			iter.Release()
			return
		}
	}
	iter.Release()
}

// 遍历表数据，执行函数为参数
func (s *Select) ForRDFun(asc bool, f func(k, v []byte) bool) {
	//s.FindPrefixFun([]byte(s.Tbname+Split), asc, f)
	s.FindPrefixFun(s.GetTbKey(), asc, f)
}

// 获取表数据
func (s *Select) ForRD(asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) {
	r = s.FindPrefix(s.GetTbKey(), asc, b, count, []int{}, distinct)
	return
}

// 统计表的记录数
func (s *Select) Count() (r int) {
	iter, ok := s.IterPrefixMove(s.GetTbKey(), true)
	if !ok {
		return
	}
	r = NewIters(iter, ok, true, 0, -1).ForDataCount()
	return
}

// 遍历表数据集
func (s *Select) For(f func(k, v []byte) bool) {
	//s.FindPrefixFun([]byte(s.Tbname), true, f)
	s.FindPrefixFun(s.GetTbKey(), true, f)
}

/*
前缀遍历
bint 第几条开始
asc,升/降序
*/
func (s *Select) FindPrefix(key []byte, asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) {
	iter, ok := s.IterPrefixMove(key, asc)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForData(showFileds, distinct)
	return
}

// 前缀遍历,执行函数为参数
func (s *Select) FindPrefixFun(key []byte, asc bool, f func(k, v []byte) bool) {
	iter, ok := s.IterPrefixMove(key, asc)
	if !ok {
		return
	}
	NewIters(iter, ok, asc, 0, -1).ForKVFun(f)
}

// 前缀遍历,统计记录数
func (s *Select) WhereIdxCount(fieldname, fieldvalue []byte) (r int) {
	key := s.GetIdxPrefix(fieldname, fieldvalue)
	iter, ok := s.IterPrefixMove(key, true)
	if !ok {
		return
	}
	r = NewIters(iter, ok, true, 0, -1).ForDataCount()
	return
}

// 根据前缀判断是否存在数据
func (s *Select) WhereIdxExist(fieldname, fieldvalue []byte) (r bool) {
	key := s.GetIdxPrefix(fieldname, fieldvalue)
	_, r = s.IterPrefixMove(key, true)
	return
}

/*
范围遍历
asc,升/降序
*/
func (s *Select) FindRand(bkey, ekey []byte, asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) {
	iter, ok := s.IterRandMove(bkey, ekey, asc)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForData(showFileds, distinct)
	return
}

// 范围遍历,执行函数为参数
func (s *Select) FindRandFun(bkey, ekey []byte, asc bool, f func(k, v []byte) bool) {
	iter, ok := s.IterRandMove(bkey, ekey, asc)
	if !ok {
		return
	}
	NewIters(iter, ok, asc, 0, -1).ForKVFun(f)
}

/*
定位遍历
asc,升/降序
*/
func (s *Select) FindSeek(key []byte, asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) {
	iter, ok := s.IterSeekMove(key)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForData(showFileds, distinct)
	return
}

// 定位遍历,执行函数为参数
func (s *Select) FindSeekFun(key []byte, asc bool, f func(k, v []byte) bool) {
	iter, ok := s.IterSeekMove(key)
	if !ok {
		return
	}
	NewIters(iter, ok, asc, 0, -1).ForKVFun(f)
}

// 获取一个key的values
func (s *Select) GetValue(key []byte) (r []byte) {
	r, _ = s.db.Get(key, nil)
	return
}

// 根据主键值获取表的一条记录value（获取一个key的value）
func (s *Select) GetPKValue(fieldvalue []byte) (r []byte) { //GetRecord
	key := s.GetPkKey(fieldvalue)
	if key == nil {
		return
	}
	r = s.GetValue(key)
	return
}

// 根据主键获取表的一条记录（获取一个key的values）
func (s *Select) Record(PKvalue []byte, showFileds []int) (r *TbData) { //GetOneRecord
	key := s.GetPkKey(PKvalue)
	value := s.GetValue(key)
	if len(value) == 0 {
		return
	}
	//r = TbDatapool.Get().(*TbData)
	//r.Release() //确保数据不混乱
	r = TbDatapool.New().(*TbData)
	r.Rd = append(r.Rd, KVToRd(key, value, showFileds))
	return

}

// 根据多个主键获取表对应的多条记录
func (s *Select) Records(PKids [][]byte, showFileds []int) (r *TbData) {
	var value []byte
	//r = TbDatapool.Get().(*TbData)
	//r.Release()
	r = TbDatapool.New().(*TbData)
	var pk []byte
	for _, v := range PKids {
		//value = s.GetValue(v)
		pk = s.GetPkKey(v)
		value = s.GetValue(pk)
		r.Rd = append(r.Rd, KVToRd(pk, value, showFileds))
	}
	return
}

// 根据主键区间获取表的区间记录
func (s *Select) RecordRand(bpk, epk []byte, showFileds []int, distinct bool) (r *TbData) {
	bid := s.GetPkKey(bpk) //t.Ifo.FieldChByte(t.Ifo.Fields[0], bpk)
	eid := s.GetPkKey(epk) //t.Ifo.FieldChByte(t.Ifo.Fields[0], epk)
	r = s.FindRand(bid, eid, true, 0, -1, showFileds, distinct)
	return
}

// 根据索引记录列表返回表记录数据
// b，开始记录，count，返回条数
func (s *Select) WhereIdx(fieldname, value []byte, asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) { //GetTableRecordForIdx
	r = s.WhereIdxs(fieldname, value, asc, b, count, showFileds, true, distinct)
	return
}

// 根据索引匹配记录列表返回表记录数据，相当于sql的like语句
// b，开始记录，count，返回条数
func (s *Select) WhereIdxLike(fieldname, value []byte, asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) { //GetTableRecordForIdx
	r = s.WhereIdxs(fieldname, value, asc, b, count, showFileds, false, distinct)
	return
}

// 根据索引等于或匹配记录列表返回表记录数据
// b，开始记录，count，返回条数
func (s *Select) WhereIdxs(fieldname, value []byte, asc bool, b, count int, showFileds []int, eq bool, distinct bool) (r *TbData) { //GetTableRecordForIdx
	gip := map[bool]func(fn, fv []byte) []byte{
		true:  s.GetIdxPrefix,
		false: s.GetIdxPrefixLike,
	}
	key := gip[eq](fieldname, value) //ca.fid-
	//获取索引对应的主键集
	//索引的情况下，倒数第二个是主键值,k的倒数第一个，由于加上v（一般都是空值。全文索引则是位置），则是倒数第二个。
	//在主键情况下，则是[]int{0}
	tbd := s.FindPrefix(key, asc, b, count, []int{-2}, false) //索引是不会存在重复的。故而distinct=false。以免空耗性能。
	if tbd == nil {
		return
	}
	r = s.IdxsGetRecords(tbd, showFileds, distinct)
	tbd.Release()
	return
}

// 根据组合索引记录列表返回表记录数据
// b，开始记录，count，返回条数
func (s *Select) WhereJionIdx(idxfields, idxvalue [][]byte, value []byte, asc bool, b, count int, showFileds []int, distinct bool) (r *TbData) { //GetTableRecordForIdx
	jionkey := s.GetIdxsPrefixKey(idxfields, idxvalue, value)
	idtbd := s.FindPrefix(jionkey, asc, b, count, []int{-2}, distinct) //[]int{-2}，获取主键值
	if idtbd == nil {
		return
	}
	r = s.IdxsGetRecords(idtbd, showFileds, distinct)
	idtbd.Release()
	return
}

/*
func (s *Select) WhereJionIdx(idxfields, idxvalue [][]byte, value []byte, b, count int) (r *TbData) {
	s.WhereJionIdxs(idxfields, idxvalue, []byte{}, false, b, count, false)
	return
}
*/
// 根据索引集合获取对应的记录集合
func (s *Select) IdxsGetRecords(tbd *TbData, showFileds []int, distinct bool) (r *TbData) {
	var pkval []byte
	//r = TbDatapool.Get().(*TbData)
	//同一个函数内不能同时使用2个get。
	//因为tbd此时未释放，r = TbDatapool.Get()会将tbd偷走。
	//这里需要用New()
	r = TbDatapool.New().(*TbData)
	var pk []byte
	for _, v := range tbd.Rd {
		/*
			ks := bytes.Split(v, []byte(Split))
			//可查看KVToRD()之后rd结构
			k = ks[len(ks)-2] //单个索引和组合索引2种情况下都是倒数第2个。
			//k = ks[1]，单个索引下正确。
			//k = ks[0]主键下，即是第一个。
		*/
		pk = s.GetPkKey(v)
		pkval = s.GetValue(pk) //s.GetPKValue(k)
		if pkval != nil {      //KVToRd(v, value, showFileds)
			rd := KVToRd(pk, pkval, showFileds)
			if distinct {
				if ArraysContains(r.Rd, rd) {
					continue
				}
			}
			r.Rd = append(r.Rd, rd)
		}
	}
	return
}

/*
//根据索引集合获取对应的记录集合
func (s *Select) IdxsGetRecords(tbd *TbData) (r *TbData) {
	var pkval, k []byte
	//r = TbDatapool.Get().(*TbData)
	//同一个函数内不能同时使用2个get。
	//因为tbd此时未释放，r = TbDatapool.Get()会将tbd偷走。
	//这里需要用New()
	r = TbDatapool.New().(*TbData)
	for _, v := range tbd.Rd {
		ks := bytes.Split(v, []byte(Split))
		//可查看KVToRD()之后rd结构
		//k = ks[len(ks)-2]
		k = ks[1]
		pkval = s.GetPKValue(k)
		if pkval != nil {
			r.Rd = append(r.Rd, JoinBytes(k, []byte(Split), pkval))
		}
	}
	return
}
*/
//根据索引查询表记录，执行函数为参数
func (s *Select) WhereIdxFun(fieldname, value []byte, asc bool, distinct bool, f func(rd []byte) bool) {
	key := s.GetIdxPrefix(fieldname, value)
	tbd := s.FindPrefix(key, asc, 0, -1, []int{}, distinct)
	if tbd == nil {
		return
	}
	for _, v := range tbd.Rd {
		if !f(v) {
			return
		}
	}
	tbd.Release()
}

// 根据索引匹配查询表记录，执行函数为参数
func (s *Select) WhereIdxLikeFun(fieldname, value []byte, asc bool, f func(k, v []byte) bool) {
	key := s.GetIdxPrefixLike(fieldname, value)
	s.FindPrefixFun(key, asc, f)
}

/*
//根据根据主键值获取一条数据
func (s *Select) WherePK(value []byte) (r *TbData) { //GetTableRecordForIdx
	//key := s.GetPkKey(value)
	r = s.Record(value)
	return
}
*/
//根据根据主键值匹配获取数据，仅主键为字符串时有效
//b，开始记录，count，返回条数
func (s *Select) WherePKLike(value []byte, asc bool, b, count int, showFileds []int) (r *TbData) { //GetTableRecordForIdx
	key := s.GetPkKey(value)
	tbd := s.FindPrefix(key, asc, b, count, showFileds, false) //主键是不会存在重复的。故而distinct=false。以免空耗性能。
	if tbd == nil {
		return
	}
	r = s.IdxsGetRecords(tbd, showFileds, false)
	//r = s.OneRecord(key)
	return
}

// 根据根据主键值匹配获取数据，仅主键为字符串时有效。执行函数为参数
// b，开始记录，count，返回条数
func (s *Select) WherePKLikeFun(value []byte, b, count int, asc bool, f func(k, v []byte) bool) { //GetTableRecordForIdx
	key := s.GetPkKey(value)
	s.FindPrefixFun(key, asc, f)
}

// 根据根据主键范围值获取数据
// b，开始记录，count，返回条数
func (s *Select) WherePKRand(minvalue, maxvalue []byte, asc bool, b int, count int, showFileds []int, distinct bool) (r *TbData) { //GetTableRecordForIdx
	minkey := s.GetPkKey(minvalue)
	maxkey := s.GetPkKey(maxvalue)
	if string(maxkey) < string(minkey) {
		return
	}
	r = s.FindRand(minkey, maxkey, asc, b, count, showFileds, distinct)
	return
}

// 根据根据主键范围值获取数据,执行函数为参数
func (s *Select) WherePKRandFun(minvalue, maxvalue []byte, asc bool, f func(k, v []byte) bool) { //GetTableRecordForIdx
	minkey := s.GetPkKey(minvalue)
	maxkey := s.GetPkKey(maxvalue)
	if string(maxkey) < string(minkey) {
		return
	}
	s.FindRandFun(minkey, maxkey, asc, f)
}
