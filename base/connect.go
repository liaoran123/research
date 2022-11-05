package base

import (
	"research/pubgo"
	"strconv"
)

var Con *Connect

//数据库连接管理
type Connect struct {
	count       int      //索引分库个数
	fulltext    *Levdb   //文章数据库
	fulltextIdx []*Levdb //对应索引分库连接
}

func NewConnect() *Connect {
	//建立数据库连接
	count := 0
	if ConfigMap["dbidxcount"] != nil {
		count = int(ConfigMap["dbidxcount"].(float64))
	}
	path := pubgo.GetCurrentAbPath()
	fulltext := Newleveldbs(path + "db")
	var fulltextIdx []*Levdb
	for i := 0; i < count; i++ {
		fulltextIdx = append(fulltextIdx, Newleveldbs(path+"dbidx"+strconv.Itoa(i)))
	}
	return &Connect{
		count:       count,
		fulltext:    fulltext,
		fulltextIdx: fulltextIdx,
	}
}

/*
//建立数据库连接
func (c *Connect) CreateConn() {
	c.fulltext = Newleveldbs("db")
	for i := 0; i < c.count; i++ {
		c.fulltextIdx = append(c.fulltextIdx, Newleveldbs("dbidx"+strconv.Itoa(i)))
	}
}
*/

//获取索引库，由于兼容企业版分库，故而稍微麻烦。
func (c *Connect) Getidxdb(word string) *Levdb {
	if c.count > 0 { //是企业版则索引分库
		return c.fulltextIdx[int([]rune(word)[0])%c.count]
	} else {
		return c.Getartdb()
	}
}

//文章的库连接
func (c *Connect) Getartdb() *Levdb {
	return c.fulltext
}
