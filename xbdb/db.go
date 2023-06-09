// 小白数据库
// 表信息
package xbdb

import (
	"fmt"
	"log"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Xb struct {
	Db *leveldb.DB
}

// 创建或打开数据库
func NewDb(dbpath string) *Xb {
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &Xb{
		Db: db,
	}
}

// 创建所有表的操作结构
// 所有表的操作都通过*Table进行操作。
func (d *Xb) GetTables() (Tables map[string]*Table) {
	Tbnames := d.GetTbnames()
	Tables = make(map[string]*Table)
	for _, v := range Tbnames {
		Tables[v] = NewTable(d.Db, v)
	}
	return
}

/*
func (d *Xb) SetTables() { //OpenTableStructs() {
	d.SetTbnames()
	d.Tables = make(map[string]*Table)
	for _, v := range d.Tbnames {
		d.Tables[v] = NewTable(d.Db, v)
	}
}
*/
// 获取数据库所有的表名称
func (d *Xb) GetTbnames() (Tbnames []string) {
	iter := d.Db.NewIterator(util.BytesPrefix([]byte(Tbspfx+Split)), nil)
	tbname := ""
	for iter.Next() {
		tbname = strings.Split(string(iter.Key()), Split)[1]
		Tbnames = append(Tbnames, tbname)
	}
	iter.Release()
	if iter.Error() != nil {
		fmt.Printf("iter.Error(): %v\n", iter.Error())
	}
	return
}

/*
func (d *Xb) SetTbnames() {
	iter := d.Db.NewIterator(util.BytesPrefix([]byte(Tbspfx+Split)), nil)
	tbname := ""
	//var Tbnames []string
	for iter.Next() {
		tbname = strings.Split(string(iter.Key()), Split)[1]
		d.Tbnames = append(d.Tbnames, tbname)
	}
	iter.Release()
	if iter.Error() != nil {
		fmt.Printf("iter.Error(): %v\n", iter.Error())
	}
	//d.Tbnames = Tbnames

}*/
