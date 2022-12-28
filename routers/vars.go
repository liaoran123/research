package routers

import (
	"research/xbdb"
	"sync"
)

var (
	ConfigMap map[string]interface{} //配置文件
	mu        sync.RWMutex
	tifo      *xbdb.TableInfo
)

/*
//目录信息
type CataInfo struct {
	Id     int    `json:"Id"`     //目录id
	Fid    int    `json:"Fid"`    //目录fid
	Title  string `json:"Title"`  //目录名称
	Isleaf string `json:"Isleaf"` //是否叶子目录

}

//所有子目录信息
type CataInfos struct {
	Catas []CataInfo `json:"Catainfo"`
}
*/
