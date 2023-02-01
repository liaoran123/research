//小白数据库

package xbdb

import (
	"bytes"
	"fmt"
)

//sql的Group语句
//支持多个字段sum
type Group struct {
	//数组对应多个字段sum
	floatvalue []map[string]float64 //float64兼容int, 和 int64和float32
}

func NewGroup() *Group {
	return &Group{}
}
func (s *Group) Do(Tbd *TbData) *TbData {
	tdb := TbDatapool.Get().(*TbData)
	tdb.Release() //确保如果上次忘记释放内存不会导致数据混乱。
	s.sum(Tbd)
	mlen := len(s.floatvalue)
	var avg float64
	for k, v := range s.floatvalue[0] {
		var trd []byte
		trd = JoinBytes([]byte(k), []byte(Split), SplitToCh(Float64ToByte(v)), []byte(Split), SplitToCh(Float64ToByte(v/float64(mlen))))
		fmt.Println(k, v, v/float64(mlen)) //分组名称，sum、avg
		for i := 1; i < mlen; i++ {        //多个sum、avg情况
			fmt.Println(s.floatvalue[i][k] / float64(mlen))
			avg = s.floatvalue[i][k] / float64(mlen)
			trd = JoinBytes(trd, []byte(Split), SplitToCh(Float64ToByte(avg)))
		}
		tdb.Rd = append(tdb.Rd, trd)
	}
	return tdb
}

//Tbd.Rd第一个字段是要分组的字段，后面有多少个字段就是sum多少个
func (s *Group) sum(Tbd *TbData) {
	var ks [][]byte
	sumlen := len(bytes.Split(Tbd.Rd[0], []byte(Split))) - 1
	for i := 0; i < sumlen; i++ {
		s.floatvalue = append(s.floatvalue, nil)
	}
	for _, v := range Tbd.Rd {
		ks = bytes.Split(v, []byte(Split))
		for i := 1; i < sumlen; i++ { //ks[0]是Group的字段，后面的都是要sum的字段，可以是多个
			s.sumfloat64(i, string(ks[0]), ks[i])
		}
	}
}

func (s *Group) sumfloat64(i int, fieldname string, fieldvalue []byte) {
	if v, ok := s.floatvalue[i][fieldname]; ok {
		s.floatvalue[i][fieldname] = v + ByteToFloat64(fieldvalue)
	} else {
		s.floatvalue[i][fieldname] = ByteToFloat64(fieldvalue)
	}
}
