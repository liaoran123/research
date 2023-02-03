//小白数据库
//表信息
package xbdb

//修改整条记录。等于先删除后添加。
func (t *Table) Updata(vals [][]byte) (r ReInfo) {
	r = t.Delete(vals[0])
	if !r.Succ {
		return
	}
	r = t.Insert(vals)
	return
}

//修改某个某些字段数据。不需要修改的字段索引不会删除和重复添加，提高性能。
func (t *Table) Upd(params map[string]string) (r ReInfo) {
	var updatefield []bool
	uvals := t.StrToByte(params)
	for _, v := range uvals {
		if len(v) == 0 {
			updatefield = append(updatefield, false)
		} else {
			updatefield = append(updatefield, true)
		}
	}
	key := JoinBytes(t.Select.GetTbKey(), uvals[0])
	data, err := t.db.Get(key, nil) //获取旧数据
	if err != nil {
		r.Info = err.Error()
		return
	}
	//组织数据
	var dvals [][]byte
	dvals = append(dvals, uvals[0])
	dvals = append(dvals, SplitRd(data)...)
	r = t.Acts(dvals, "delete", updatefield) //删除旧数据
	if !r.Succ {
		return
	}
	//更新数据
	for i, v := range uvals {
		if len(v) != 0 { //即是要修改的字段
			dvals[i] = v //更改要更新的字段值
		}
	}
	r = t.Acts(dvals, "insert", updatefield) //添加新数据
	return
}
