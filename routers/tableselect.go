package routers

import (
	"bytes"
	"research/xbdb"
	"strconv"
	"strings"
	"sync"
)

var bufpool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

/*
//根据索引查询表记录并组织json数据
func IdxGetDataToJson(tbname, idxfield, idxvalue []byte) (r *bytes.Buffer) {
	kvs := Table[string(tbname)].Select.WhereIdx(idxfield, idxvalue, 0, -1)
	if kvs == nil {
		return
	}
	r = DataToJson(kvs, Table[string(tbname)].Ifo)
	return
}
*/
//根据表记录组织json数据
func DataToJson(tbd *xbdb.TbData, Seltbifo xbdb.TableInfo) (r *bytes.Buffer) {
	if tbd == nil {
		return
	}
	r = bufpool.Get().(*bytes.Buffer)
	var value [][]byte
	jsonstr := ""
	/*
		[{"id":2,"title":"金刚经","fid":1,"isleaf":"0"},
		{"id":3,"title":"六祖坛经","fid":1,"isleaf":"0"}]
	*/
	r.WriteString("[")
	for j, v := range tbd.Rd {
		if v == nil {
			continue
		}
		r.WriteString("{")
		value = bytes.Split(v, []byte(xbdb.Split))
		for i, fv := range Seltbifo.FieldType {
			switch fv {
			case "string":
				jsonstr = "\"" + Seltbifo.Fields[i] + "\":" + strconv.Quote(string(value[i])) //strconv.Quote自动加字符串号
			default:
				iv := Seltbifo.ByteChString(Seltbifo.FieldType[i], value[i])
				jsonstr = "\"" + Seltbifo.Fields[i] + "\":" + iv
			}
			if i != len(Seltbifo.FieldType)-1 {
				jsonstr += ","
			}
			jsonstr = strings.Replace(jsonstr, "\n", "\\n", -1) //json转义
			/*
				jsonstr = strings.Replace(jsonstr, "\t", "\\t", -1) //json转义
				jsonstr = strings.Replace(jsonstr, "\n", "\\n", -1) //json转义
								content = strings.Replace(content, "\\u003c", "<", -1)
					content = strings.Replace(content, "\\u003e", ">", -1)
					content = strings.Replace(content, "\\u0026", "&", -1)
			*/
			r.WriteString(jsonstr)
		}
		r.WriteString("}")
		if j != len(tbd.Rd)-1 {
			r.WriteString(",")
		}
	}
	r.WriteString("]")
	tbd.Release()
	return
}
