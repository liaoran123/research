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
	var vals [][]byte
	f, v := "", ""
	for i := 0; i < len(Table[name].Ifo.FieldType); i++ {
		f = Table[name].Ifo.FieldType[i]
		v = Table[name].Ifo.Fields[i]
		vals = append(vals, Table[name].Ifo.TypeChByte(f, params[v]))
	}
	if insorupd == "" {
		insorupd = "ins"
	}
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

//删除一条表记录
func delete(tbname, id string) (r xbdb.ReInfo) {
	bid := Table[tbname].Ifo.TypeChByte(Table[tbname].Ifo.FieldType[0], id)
	r = Table[tbname].Delete(bid)
	return
}

//获取表的一条记录
func getonerecord(tbname, idxfield, id string, w http.ResponseWriter) {
	idxvalue := Table[tbname].Ifo.FieldChByte(idxfield, id)
	tbd := Table[tbname].Select.OneRecord(idxvalue)
	if tbd == nil {
		return
	}
	r := DataToJson(tbd, Table[tbname].Ifo)
	w.Write(r.Bytes())
	r.Reset()
	bufpool.Put(r)
}
