package xbdb

//通过getprefix.go表的前缀规则获取的key，获取Iter游标数据
import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//空游标，整个数据库数据游标
func (s *Select) Nil() (iter iterator.Iterator) {
	iter = s.db.NewIterator(nil, nil)
	return
}

//前缀匹配数据游标
func (s *Select) IterPrefix(key []byte) (iter iterator.Iterator) {
	iter = s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)
	return
}

//范围数据游标
func (s *Select) IterRand(b, e []byte) (iter iterator.Iterator) {
	iter = s.db.NewIterator(&util.Range{Start: b, Limit: e}, nil)
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
