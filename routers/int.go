package routers

import (
	"fmt"
	"log"
	"research/xbdb"

	"github.com/syndtr/goleveldb/leveldb"
)

var Xb *leveldb.DB
var Table map[string]*xbdb.Table

func Ini() {
	//打开或创建数据库
	dbpath := ConfigMap["dbpath"].(string)
	xb, err := leveldb.OpenFile(dbpath+"db", nil)
	if err != nil {
		log.Fatal(err)
	}
	//建表
	Xb = xb
	dbinfo := xbdb.NewTableInfo(Xb)
	if dbinfo.GetInfo("ca").FieldType == nil {
		createca(dbinfo)
	}
	if dbinfo.GetInfo("art").FieldType == nil {
		createart(dbinfo)
	}
	if dbinfo.GetInfo("c").FieldType == nil {
		createc(dbinfo)
	}
	//打开表操作结构
	Table = make(map[string]*xbdb.Table)
	Table["ca"] = xbdb.NewTable(Xb, "ca")
	Table["art"] = xbdb.NewTable(Xb, "art")
	Table["c"] = xbdb.NewTable(Xb, "c")
	//目录入加载内存
	CRAMs = NewCataRAMs()
	CRAMs.LoadCataRAM()
	//文章对应的目录fid加载入内存
	LoadartRAM()
}

//创建目录表
func createca(tbifo *xbdb.TableInfo) {
	name := "ca"                                            //目录表
	fields := []string{"id", "title", "fid", "isleaf"}      //字段，编码，标题，父id，是否叶子目录。1是0否。
	fieldType := []string{"int", "string", "int", "string"} //字段
	idxs := []string{"2"}                                   //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{}                                  //考据级全文搜索索引字段的下标。
	ftlen := "7"                                            //全文搜索的长度，中文默认是7
	r := tbifo.Create(name, ftlen, fields, fieldType, idxs, fullText)
	fmt.Printf("r: %v\n", r)
}

//创建文章表
func createart(tbifo *xbdb.TableInfo) {
	name := "art"                                                     //目录表
	fields := []string{"id", "title", "fid", "split", "url"}          //字段，编码，标题，目录/父id，割截内容符号，网址
	fieldType := []string{"int", "string", "int", "string", "string"} //字段
	idxs := []string{"2"}                                             ////索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{}                                            //考据级全文搜索索引字段的下标。
	ftlen := "7"                                                      //全文搜索的长度，中文默认是7
	r := tbifo.Create(name, ftlen, fields, fieldType, idxs, fullText)
	fmt.Printf("r: %v\n", r)
}

//创建文章内容表，该表是全文搜索，故而名称尽量短，可以减少文件大小。
//带全文搜索索引的内容表c
func createc(tbifo *xbdb.TableInfo) {
	name := "c"                               //目录表，
	fields := []string{"id", "s"}             //字段 该id:=art.id,secid为字符串。s 是文章的分段内容,pos,为位置
	fieldType := []string{"string", "string"} //字段
	idxs := []string{}                        //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{"1"}                 //考据级全文搜索索引字段的下标。
	ftlen := "7"                              //全文搜索的长度，中文默认是7
	r := tbifo.Create(name, ftlen, fields, fieldType, idxs, fullText)
	fmt.Printf("r: %v\n", r)
}
