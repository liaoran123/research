//小白数据库

package xbdb

import "bytes"

//查询要显示的字段。如SQL语句：select f0,f1,f2 from tb 。f0,f1,f2即是要显示的字段
type Selectfiled struct {
	fieldidxs []int
	Tbd       *TbData
}

func NewSelectfiled() *Selectfiled {
	return &Selectfiled{
		Tbd: TbDatapool.Get().(*TbData),
	}
}

//重复用Tbd
func (s *Selectfiled) NewTbd() {
	s.Tbd.Release()
	s.Tbd = TbDatapool.Get().(*TbData)
}

//获取要显示的字段idx
func (s *Selectfiled) Setfieldidxs(Ifo *TableInfo, fieldnames []string) {
	if len(s.fieldidxs) != 0 {
		s.fieldidxs = s.fieldidxs[:0]
	}
	//查找对应的字段的idx
	for _, fv := range fieldnames {
		for i, tv := range Ifo.Fields {
			if fv == tv {
				s.fieldidxs = append(s.fieldidxs, i)
				break
			}
		}
	}
}

//过滤掉不需要显示的字段
func (s *Selectfiled) Filter(rd []byte) bool {
	ks := bytes.Split(rd, []byte(Split))
	var frd []byte
	for _, v := range s.fieldidxs {
		frd = JoinBytes(frd, ks[v])
	}
	s.Tbd.Rd = append(s.Tbd.Rd, frd)
	return true
}
