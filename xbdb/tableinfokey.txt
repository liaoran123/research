//小白数据库
//获取各种key值
//该文件已经转移到Select，可以删除
package xbdb

/*
//一条索引key
func (t *TableInfo) getIdxPfxKey000(idxfield, idxvalue, pkvalue []byte) (r []byte) {
	//bSplit := []byte(Split)
	//r = JoinBytes([]byte(t.Name+","), idxfield, bSplit, idxvalue, bSplit, pkvalue)
	r = GetIdxPfxKey([]byte(t.Name), idxfield, idxvalue, pkvalue)
	return
}
*/
//一条主键key
func GetPkKey(tbname, pkvalue []byte) (r []byte) {
	bSplit := []byte(Split)
	r = JoinBytes([]byte(tbname), bSplit, pkvalue)
	return
}

//一条索引key
func GetIdxPfxKey(tbname, idxfield, idxvalue, pkvalue []byte) (r []byte) {
	bSplit := []byte(Split)
	r = JoinBytes(tbname, []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit, pkvalue)
	return
}

//一条组合索引key
func GetIdxsPfxKey(tbname, pkvalue []byte, idxfields, idxvalues [][]byte) (r []byte) {
	bSplit := []byte(Split)
	flen, ilen := len(idxfields), len(idxvalues)
	if flen != ilen {
		return
	}
	var idxfield, idxvalue []byte
	for i := 0; i < ilen; i++ {
		idxfield = JoinBytes(idxfields[i])
		idxvalue = JoinBytes(idxvalues[i])
		if i != ilen-1 {
			idxfield = JoinBytes(idxfield, []byte(IdxSplit))
			idxvalue = JoinBytes(idxvalue, []byte(IdxSplit))
		}
	}
	r = JoinBytes(tbname, []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit, pkvalue)
	return
}

//索引前缀
func GetIdxPfx(tbname, idxfield, idxvalue []byte) (r []byte) {
	bSplit := []byte(Split)
	r = JoinBytes(tbname, []byte(IdxSplit), idxfield, bSplit, idxvalue, bSplit)
	return
}
