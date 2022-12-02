package base

import (
	"strings"

	"github.com/syndtr/goleveldb/leveldb/util"
)

//目录fid索引
type catafididx struct {
	tbn string
}

func Newcatafididx() *catafididx {
	return &catafididx{
		tbn: "cf", //表名

	}
}

//k=cf-fid-cataid
func (c *catafididx) setkey(fid, cataid int) (r []byte) {
	r = JoinBytes([]byte(c.tbn+"~"), IntToBytes(fid), []byte("~"), IntToBytes(cataid))
	return
}

//添加目录fid索引
func (c *catafididx) Insert(fid, cataid int) (r bool) {
	err = Con.Getartdb().Db.Put(c.setkey(fid, cataid), []byte{}, nil)
	//err = Con.Getartdb().Db.Put(JoinBytes([]byte(c.tbn+"~"), IntToBytes(fid), []byte("~"), IntToBytes(cataid)), []byte{}, nil) //添加目录标题
	Chekerr()
	r = err == nil
	return
}

//删除目录fid索引
func (c *catafididx) Delete(fid, cataid int) (r bool) {
	err = Con.Getartdb().Db.Delete(c.setkey(fid, cataid), nil)
	//err = Con.Getartdb().Db.Put(JoinBytes([]byte(c.tbn+"~"), IntToBytes(fid), []byte("~"), IntToBytes(cataid)), []byte{}, nil) //添加目录标题
	Chekerr()
	r = err == nil
	return
}

//获取所有子目录
func (c *catafididx) ChildCatas(fid int) (r []int) {
	iter := Con.Getartdb().Db.NewIterator(util.BytesPrefix(JoinBytes([]byte(c.tbn+"~"), IntToBytes(fid), []byte("~"))), nil)
	var ks []string
	var cataid int
	for iter.Next() {
		ks = strings.Split(string(iter.Key()), "~")
		cataid = BytesToInt([]byte(ks[2]))
		r = append(r, cataid)
	}
	Release(iter)
	return
}
