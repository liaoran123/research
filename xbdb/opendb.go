//小白数据库
//表信息
package xbdb

import (
	"fmt"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var xb *leveldb.DB

//创建或打开数据库
func OpenDb(fp string) error {
	xb, err = leveldb.OpenFile(fp+"db", nil)
	return err
}

//创建所有表的操作结构
func OpenTables() map[string]*Table {
	iter := xb.NewIterator(util.BytesPrefix([]byte(Tbspfx+Split)), nil)
	Tables := make(map[string]*Table)
	tbname := ""
	for iter.Next() {
		tbname = strings.Split(string(iter.Key()), Split)[1]
		Tables[tbname] = NewTable(tbname)
	}
	iter.Release()
	if iter.Error() != nil {
		fmt.Printf("iter.Error(): %v\n", iter.Error())
	}
	return Tables
}
