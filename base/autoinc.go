package base

import (
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb/util"
)

//目录和网站各一个公共自动增值
var PcataAutoinc, PcontenAutoinc *Autoinc

//自动增值
type Autoinc struct {
	tbn string
	id  int
	//ld  *Levdb //文章数据库
	mu sync.RWMutex
}

func NewAutoinc(tbn string) *Autoinc {
	var id int
	//查询该表最大值+1，为当前自动增值的值。
	iter := Con.Getartdb().Db.NewIterator(util.BytesPrefix([]byte(tbn+"-")), nil)
	if iter.Last() {
		lid := strings.Split(string(iter.Key()), "-")[1]
		id = BytesToInt([]byte(lid)) + 1
	}
	if id == 0 {
		id = 1 //id从1开始。当fid=0，则能表示为顶级目录，即没有父id
	}
	return &Autoinc{
		id: id,
		//ld:  Con.Getartdb(),
		tbn: tbn,
	}
}

func (a *Autoinc) Getid() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	id := a.id
	a.id++
	return id
}

/*
//将自动增值最后id写入数据库表
func (a *Autoinc) Writelastid() (r bool) {
	a.mu.Lock()
	err = Con.Getartdb().Db.Put([]byte(a.tbn+"lastid"), IntToBytes(a.id), nil)
	a.mu.Unlock()
	Chekerr()
	r = err == nil
	return
}
*/
/*
func (a *Autoinc) Setid(id int) {
	a.id = id
}
//+1
func (a *Autoinc) Inc() {
	a.id++
}
*/
