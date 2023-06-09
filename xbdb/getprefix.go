package xbdb

//表的前缀规则
import (
	"bytes"
)

// 一条表前缀，可以遍历表所有。包括记录和索引等等。
func (s *Select) GetTbLikeKey() (r []byte) {
	r = []byte(s.Tbname)
	return
}

// 一条表前缀，可以遍历表的记录。
func (s *Select) GetTbKey() (r []byte) {
	r = []byte(s.Tbname + Split)
	return
}

// 一条主键key
func (s *Select) GetPkKey(pkvalue []byte) (r []byte) {
	//bSplit := []byte(Split)
	//r = JoinBytes([]byte(s.Tbname), bSplit, pkvalue)
	r = JoinBytes(s.GetTbKey(), pkvalue)
	return
}

/*
//一条主键key，前缀匹配，仅当主键为字符串时有效
func (s *Select) GetPkKeyLike(pkvalue []byte) (r []byte) {
	r = s.GetPkKey(pkvalue)
	r = bytes.Trim(r, Split)
	return
}
*/
//一条索引key
func (s *Select) GetIdxPrefixKey(idxfield, idxvalue, pkvalue []byte) (r []byte) {
	bSplit := []byte(Split)
	r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit, pkvalue)
	return
}

// 一条组合索引key，GetIdxPrefixKey是一个值，GetIdxsPrefixKey是的多个值
func (s *Select) GetIdxsPrefixKey(idxfield, idxvalue [][]byte, pkvalue []byte) (r []byte) {
	bIdxSplit := []byte(IdxSplit)                //索引拼接分隔符
	idxfields := bytes.Join(idxfield, bIdxSplit) //只需将多个值拼接起来即可
	idxvalues := bytes.Join(idxvalue, bIdxSplit) //只需将多个值拼接起来即可
	r = s.GetIdxPrefixKey(idxfields, idxvalues, pkvalue)
	//bSplit := []byte(Split)
	//r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfields, bSplit, idxvalues, bSplit, pkvalue)
	return
}

// 一条组合索引的like前缀,第一个或或nil值前的数据为前缀
func (s *Select) GetIdxsLikePrefixKey(idxfield, idxvalue [][]byte, pkvalue []byte) (r []byte) {
	bIdxSplit := []byte(IdxSplit)
	bSplit := []byte(Split)
	r = JoinBytesNoNil([]byte(s.Tbname), bIdxSplit, idxfield[0], bIdxSplit, idxfield[1], bSplit, idxvalue[0], bIdxSplit, idxvalue[1], bSplit, pkvalue)
	return
}

// 索引前缀，等于索引idxvalue
func (s *Select) GetIdxPrefix(idxfield, idxvalue []byte) (r []byte) {
	r = s.GetIdxPrefixKey(idxfield, idxvalue, []byte{}) //只需通过GetIdxPrefixKey，提供一个nil的pkvalue即可。
	/*
		bSplit := []byte(Split)
		r = JoinBytes([]byte(s.Tbname), []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit)
	*/
	return
}

// 索引前缀，索引idxvalue也前缀匹配。即是sql的like语句
func (s *Select) GetIdxPrefixLike(idxfield, idxvalue []byte) (r []byte) {
	r = s.GetIdxPrefix(idxfield, idxvalue)
	r = bytes.Trim(r, Split)
	return
}
