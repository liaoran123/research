package routers

import (
	"net/http"
	"research/xbdb"
)

func InsOrUpd(name string, params map[string]string, insorupd string) (r xbdb.ReInfo) {
	iou := map[string]func(name string, v [][]byte) (r xbdb.ReInfo){
		"ins": ins,
		"upd": upd,
	}
	if insorupd == "" {
		insorupd = "ins"
	}
	/*2023-01-04修改屏蔽
	var vals [][]byte
	f, v := "", ""
	for i := 0; i < len(Table[name].Ifo.FieldType); i++ {
		f = Table[name].Ifo.FieldType[i]
		v = Table[name].Ifo.Fields[i]
		vals = append(vals, Table[name].Ifo.TypeChByte(f, params[v]))
	}
	*/
	vals := Table[name].StrToByte(params) //2023-01-04增加
	r = iou[insorupd](name, vals)
	return
}
func ins(name string, vals [][]byte) (r xbdb.ReInfo) {
	r = Table[name].Insert(vals)
	return
}
func upd(name string, vals [][]byte) (r xbdb.ReInfo) {
	r = Table[name].Updata(vals)
	return
}

// 删除一条表记录
func deleteRd(tbname, id string) (r xbdb.ReInfo) {
	bid := Table[tbname].Ifo.TypeChByte(Table[tbname].Ifo.FieldType[0], id)
	r = Table[tbname].Delete(bid)
	return
}

// 获取表的一条记录
func getonerecord(tbname, idxfield, id string, w http.ResponseWriter) {
	idxvalue := Table[tbname].Ifo.FieldChByte(idxfield, id)
	tbd := Table[tbname].Select.Record(idxvalue, []int{})
	if tbd == nil {
		return
	}
	r := Table[tbname].DataToJsonApp(tbd) //r := Table[tbname].DataToJson(tbd)
	w.Write(r.Bytes())
	r.Reset()
	xbdb.Bufpool.Put(r)
}
