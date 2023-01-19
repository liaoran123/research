//小白数据库
//表信息
package xbdb

import (
	"bytes"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Select struct {
	Db     *leveldb.DB
	Tbname string
}

var (
	iterFixed map[bool]func(iter iterator.Iterator) bool //起始位置。升序，first；降序，last
	itermove  map[bool]func(iter iterator.Iterator) bool //移动netx, prev
)

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

func NewSelect(tbname string, db *leveldb.DB) *Select { //*leveldb.DB {

	iterFixed = make(map[bool]func(iter iterator.Iterator) bool, 2)
	iterFixed[true] = First
	iterFixed[false] = Last
	itermove = make(map[bool]func(iter iterator.Iterator) bool, 2)
	itermove[true] = Next
	itermove[false] = Prev

	return &Select{
		Db:     db,
		Tbname: tbname,
	}
}

//空游标，整个数据库数据游标
func (s *Select) Nil() (iter iterator.Iterator) {
	iter = s.Db.NewIterator(nil, nil)
	return
}

//前缀匹配数据游标
func (s *Select) IterPrefix(key []byte) (iter iterator.Iterator) {
	iter = s.Db.NewIterator(util.BytesPrefix([]byte(key)), nil)
	return
}

//范围数据游标
func (s *Select) IterRand(b, e []byte) (iter iterator.Iterator) {
	iter = s.Db.NewIterator(&util.Range{Start: b, Limit: e}, nil)
	return
}

//定位游标
func (s *Select) IterSeekMove(key []byte) (iter iterator.Iterator, ok bool) {
	iter = s.Nil()
	ok = iter.Seek(key)
	return
}

//前缀游标
func (s *Select) IterPrefixMove(key []byte, asc bool) (iter iterator.Iterator, ok bool) { //Prefixiter
	iter = s.IterPrefix(key)
	if iter != nil {
		ok = iterFixed[asc](iter)
	}
	return
}

//范围游标
func (s *Select) IterRandMove(bkey, ekey []byte, asc bool) (iter iterator.Iterator, ok bool) {
	iter = s.IterRand(bkey, ekey)
	if iter != nil {
		ok = iterFixed[asc](iter)
	}
	return
}

/*
遍历数据库

func (s *Select) ForDb(f func(rd []byte) bool) {
	iter := s.Nil()
	for iter.Next() {
		if f(KVToRd(iter.Key(), iter.Value())) {
		} else {
			iter.Release()
			return
		}
	}
	iter.Release()
}
*/
//遍历数据库，主要用于复制数据库
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

/*
遍历表数据

func (s *Select) ForData(f func(rd []byte) bool) {
	s.FindPrefixFun([]byte(s.Tbname+Split), true, f)
}
*/
//遍历表数据
func (s *Select) ForRDFun(asc bool, f func(rd []byte) bool) {
	s.FindPrefixFun([]byte(s.Tbname+Split), asc, f)
}

//遍历表数据
func (s *Select) ForRD(asc bool, b, count int) (r *TbData) {
	r = s.FindPrefix([]byte(s.Tbname+Split), asc, b, count)
	return
}

//统计表的记录数
func (s *Select) Count() (r int) {
	iter, ok := s.IterPrefixMove([]byte(s.Tbname+Split), true)
	if !ok {
		return
	}
	r = NewIters(iter, ok, true, 0, -1).ForDataCount()
	return
}

/*
遍历表所有
*/
func (s *Select) For(f func(rd []byte) bool) {
	s.FindPrefixFun([]byte(s.Tbname), true, f)
}

/*
前缀遍历
bint 第几条开始
asc,升/降序
*/
func (s *Select) FindPrefix(key []byte, asc bool, b, count int) (r *TbData) {
	iter, ok := s.IterPrefixMove(key, asc)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForData()
	return
}

//前缀遍历,执行函数为参数
func (s *Select) FindPrefixFun(key []byte, asc bool, f func(rd []byte) bool) {
	iter, ok := s.IterPrefixMove(key, asc)
	if !ok {
		return
	}
	NewIters(iter, ok, asc, 0, -1).ForDataFun(f)
}

//前缀遍历,统计记录数
func (s *Select) FindIDXCount(fieldname, fieldvalue []byte) (r int) {
	key := s.GetIdxPrefix(fieldname, fieldvalue)
	iter, ok := s.IterPrefixMove(key, true)
	if !ok {
		return
	}
	r = NewIters(iter, ok, true, 0, -1).ForDataCount()
	return
}

//根据前缀判断是否存在数据
func (s *Select) FindIDXExist(fieldname, fieldvalue []byte) (r bool) {
	key := s.GetIdxPrefix(fieldname, fieldvalue)
	_, r = s.IterPrefixMove(key, true)
	return
}

/*
范围遍历
asc,升/降序
*/
func (s *Select) FindRand(bkey, ekey []byte, asc bool, b, count int) (r *TbData) {
	iter, ok := s.IterRandMove(bkey, ekey, asc)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForData()
	return
}

//范围遍历,执行函数为参数
func (s *Select) FindRandFun(bkey, ekey []byte, asc bool, f func(rd []byte) bool) {
	iter, ok := s.IterRandMove(bkey, ekey, asc)
	if !ok {
		return
	}
	NewIters(iter, ok, asc, 0, -1).ForDataFun(f)
}

/*
定位遍历
asc,升/降序
*/
func (s *Select) FindSeek(key []byte, asc bool, b, count int) (r *TbData) {
	iter, ok := s.IterSeekMove(key)
	if !ok {
		return
	}
	r = NewIters(iter, ok, asc, b, count).ForData()
	return
}

//定位遍历,执行函数为参数
func (s *Select) FindSeekFun(key []byte, asc bool, f func(rd []byte) bool) {
	iter, ok := s.IterSeekMove(key)
	if !ok {
		return
	}
	NewIters(iter, ok, asc, 0, -1).ForDataFun(f)
}

//获取一个key的values
func (s *Select) GetValue(key []byte) (r []byte) {
	r, _ = s.Db.Get(key, nil)
	return
}

//一条主键key
func (s *Select) GetPkKey(pkvalue []byte) (r []byte) {
	bSplit := []byte(Split)
	r = JoinBytes([]byte(s.Tbname), bSplit, pkvalue)
	return
}

//一条主键key，前缀匹配，仅当主键为字符串时有效
func (s *Select) GetPkKeyLike(pkvalue []byte) (r []byte) {
	r = s.GetPkKey(pkvalue)
	r = bytes.Trim(r, Split)
	return
}

//根据主键值获取表的一条记录value（获取一个key的value）
func (s *Select) GetPKValue(fieldvalue []byte) (r []byte) { //GetRecord
	key := s.GetPkKey(fieldvalue)
	if key == nil {
		return
	}
	r = s.GetValue(key)
	return
}

//一条索引key
func (s *Select) GetIdxPrefixKey(idxfield, idxvalue, pkvalue []byte) (r []byte) {
	bSplit := []byte(Split)
	r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit, pkvalue)
	return
}

//一条组合索引key，GetIdxPrefixKey是一个值，GetIdxsPrefixKey是的多个值
func (s *Select) GetIdxsPrefixKey(idxfield, idxvalue [][]byte, pkvalue []byte) (r []byte) {
	bIdxSplit := []byte(IdxSplit)                //索引拼接分隔符
	idxfields := bytes.Join(idxfield, bIdxSplit) //只需将多个值拼接起来即可
	idxvalues := bytes.Join(idxvalue, bIdxSplit) //只需将多个值拼接起来即可
	r = s.GetIdxPrefixKey(idxfields, idxvalues, pkvalue)
	//bSplit := []byte(Split)
	//r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfields, bSplit, idxvalues, bSplit, pkvalue)
	return
}

/*
//一条索引key
func (s *Select) GetIdxPrefixKey(idxfield, idxvalue, pkvalue []byte) (r []byte) {
	idxfields := bytes.Split(idxfield, []byte(","))
	idxvalues := bytes.Split(idxvalue, []byte(","))
	r = s.GetIdxsPrefixKey(pkvalue, idxfields, idxvalues)
	return
}
*/

/*

//一条组合索引key，GetIdxPrefixKey是一个值，GetIdxsPrefixKey是的多个值
func (s *Select) GetIdxsPrefixKey(pkvalue []byte, idxfields, idxvalues [][]byte) (r []byte) {
	bSplit := []byte(Split)
	bIdxSplit := []byte(IdxSplit) //索引拼接分隔符
	flen, ilen := len(idxfields), len(idxvalues)
	if flen != ilen {
		return
	}
	var idxfield, idxvalue []byte
	for i := 0; i < ilen; i++ {
		idxfield = idxfields[i] //JoinBytes(idxfields[i])
		idxvalue = idxvalues[i] //JoinBytes(idxvalues[i])
		if i != ilen-1 {
			idxfield = JoinBytes(idxfield, bIdxSplit)
			idxvalue = JoinBytes(idxvalue, bIdxSplit)
		}
	}
	r = JoinBytes([]byte(s.Tbname), bIdxSplit, idxfield, bSplit, idxvalue, bSplit, pkvalue)
	return
}


//一条组合索引key
func (s *Select) GetIdxsPrefixKey(pkvalue []byte, idxfields, idxvalues [][]byte) (r []byte) {
	bSplit := []byte(Split)
	flen, ilen := len(idxfields), len(idxvalues)
	if flen != ilen {
		return
	}
	var idxfield, idxvalue []byte
	for i := 0; i < ilen; i++ {
		idxfield = idxfields[i] //JoinBytes(idxfields[i])
		idxvalue = idxvalues[i] //JoinBytes(idxvalues[i])
		if i != ilen-1 {
			idxfield = JoinBytes(idxfield, []byte(IdxSplit))
			idxvalue = JoinBytes(idxvalue, []byte(IdxSplit))
		}
	}
	r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit, pkvalue)
	return
}


//根据主键获取表的一条记录（获取一个key的values）
func (s *Select) OneRecord(PKvalue []byte) (r *TbData) { //GetOneRecord
	r = s.Record(PKvalue)
	return
}
*/

//索引前缀，等于索引idxvalue
func (s *Select) GetIdxPrefix(idxfield, idxvalue []byte) (r []byte) {
	r = s.GetIdxPrefixKey(idxfield, idxvalue, []byte{}) //只需通过GetIdxPrefixKey，提供一个nil的pkvalue即可。
	/*
		bSplit := []byte(Split)
		r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit)
	*/
	return
}

//索引前缀，索引idxvalue也前缀匹配。即是sql的like语句
func (s *Select) GetIdxPrefixLike(idxfield, idxvalue []byte) (r []byte) {
	r = s.GetIdxPrefix(idxfield, idxvalue)
	r = bytes.Trim(r, Split)
	return
}

//根据主键获取表的一条记录（获取一个key的values）
func (s *Select) Record(PKvalue []byte) (r *TbData) { //GetOneRecord
	key := s.GetPkKey(PKvalue)
	value := s.GetValue(key)
	if len(value) == 0 {
		return
	}
	r = TbDatapool.Get().(*TbData)
	r.Rd = append(r.Rd, KVToRd(key, value))
	return

}

//根据主键获取表的一条记录（获取一个key的values）
func (s *Select) Records(PKids [][]byte) (r *TbData) {
	var value []byte
	r = TbDatapool.Get().(*TbData)
	for _, v := range PKids {
		value = s.GetValue(v)
		if len(value) == 0 {
			continue
		}
		r.Rd = append(r.Rd, KVToRd(v, value))
	}
	return
}

//根据主键区间获取表的区间记录
func (s *Select) RecordRand(bpk, epk []byte) (r *TbData) {
	bid := s.GetPkKey(bpk) //t.Ifo.FieldChByte(t.Ifo.Fields[0], bpk)
	eid := s.GetPkKey(epk) //t.Ifo.FieldChByte(t.Ifo.Fields[0], epk)
	r = s.FindRand(bid, eid, true, 0, -1)
	return
}

/*
//根据索引获取索引记录列表
func (s *Select) GetRecordsForIdx(idxname, idxvalue []byte, b, count int) (r *TbData) {
	key := s.getIdxPrefix(idxname, idxvalue)
	tbd = s.FindPrefix(key, true, b, count)
	return
}
*/

//根据索引记录列表返回表记录数据
//b，开始记录，count，返回条数
func (s *Select) WhereIdx(fieldname, value []byte, asc bool, b, count int) (r *TbData) { //GetTableRecordForIdx
	r = s.WhereIdxs(fieldname, value, asc, b, count, true)
	return
}

//根据索引匹配记录列表返回表记录数据，相当于sql的like语句
//b，开始记录，count，返回条数
func (s *Select) WhereIdxLike(fieldname, value []byte, asc bool, b, count int) (r *TbData) { //GetTableRecordForIdx
	r = s.WhereIdxs(fieldname, value, asc, b, count, false)
	return
}

//根据索引等于或匹配记录列表返回表记录数据
//b，开始记录，count，返回条数
func (s *Select) WhereIdxs(fieldname, value []byte, asc bool, b, count int, eq bool) (r *TbData) { //GetTableRecordForIdx
	gip := map[bool]func(fn, fv []byte) []byte{
		true:  s.GetIdxPrefix,
		false: s.GetIdxPrefixLike,
	}
	key := gip[eq](fieldname, value) //ca.fid-
	tbd := s.FindPrefix(key, asc, b, count)
	if tbd == nil {
		return
	}
	r = s.IdxsGetRecords(tbd)
	return
}

//根据索引集合获取对应的记录集合
func (s *Select) IdxsGetRecords(tbd *TbData) (r *TbData) {
	var pkval, k []byte
	r = TbDatapool.Get().(*TbData)
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
	tbd.Release()
	return
}

//根据索引查询表记录，执行函数为参数
func (s *Select) WhereIdxFun(fieldname, value []byte, asc bool, f func(rd []byte) bool) {
	key := s.GetIdxPrefix(fieldname, value)
	tbd := s.FindPrefix(key, asc, 0, -1)
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

//根据索引匹配查询表记录，执行函数为参数
func (s *Select) WhereIdxLikeFun(fieldname, value []byte, asc bool, f func(rd []byte) bool) {
	key := s.GetIdxPrefixLike(fieldname, value)
	s.FindPrefixFun(key, asc, f)
}

//根据根据主键值获取数据
//b，开始记录，count，返回条数
func (s *Select) WherePK(value []byte) (r *TbData) { //GetTableRecordForIdx
	key := s.GetPkKey(value)
	r = s.Record(key)
	return
}

//根据根据主键值匹配获取数据，仅主键为字符串时有效
//b，开始记录，count，返回条数
func (s *Select) WherePKLike(value []byte, asc bool, b, count int) (r *TbData) { //GetTableRecordForIdx
	key := s.GetPkKeyLike(value)
	tbd := s.FindPrefix(key, asc, b, count)
	if tbd == nil {
		return
	}
	r = s.IdxsGetRecords(tbd)
	//r = s.OneRecord(key)
	return
}

//根据根据主键值匹配获取数据，仅主键为字符串时有效。执行函数为参数
//b，开始记录，count，返回条数
func (s *Select) WherePKLikeFun(value []byte, b, count int, asc bool, f func(rd []byte) bool) { //GetTableRecordForIdx
	key := s.GetPkKeyLike(value)
	s.FindPrefixFun(key, asc, f)
}

//根据根据主键范围值获取数据
//b，开始记录，count，返回条数
func (s *Select) WherePKRand(minvalue, maxvalue []byte, asc bool, b int, count int) (r *TbData) { //GetTableRecordForIdx
	minkey := s.GetPkKey(minvalue)
	maxkey := s.GetPkKey(maxvalue)
	if string(maxkey) < string(minkey) {
		return
	}
	r = s.FindRand(minkey, maxkey, asc, b, count)
	return
}

//根据根据主键范围值获取数据,执行函数为参数
func (s *Select) WherePKRandFun(minvalue, maxvalue []byte, asc bool, f func(rd []byte) bool) { //GetTableRecordForIdx
	minkey := s.GetPkKey(minvalue)
	maxkey := s.GetPkKey(maxvalue)
	if string(maxkey) < string(minkey) {
		return
	}
	s.FindRandFun(minkey, maxkey, asc, f)
}
