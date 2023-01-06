//小白数据库
//表信息
package xbdb

//修改一个数据。等于先删除后添加。
func (t *Table) Updata(vals [][]byte) (r ReInfo) {
	r = t.Delete(vals[0])
	if !r.Succ {
		return
	}
	r = t.Insert(vals)
	return
}
func (t *Table) Upd(params map[string]string) (r ReInfo) {
	vals := t.StrToByte(params)
	r = t.Updata(vals)
	return
}
