package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"research/base"
	"research/pubgo"
	"research/routers"
)

//读取配置文件装入map
func readconfig() {
	path := pubgo.GetCurrentAbPath()
	text, _ := ioutil.ReadFile(path + "config.json")
	pubgo.ConfigMap = make(map[string]interface{})
	json.Unmarshal(text, &pubgo.ConfigMap)

	text, _ = ioutil.ReadFile(path + "base/config.json")
	base.ConfigMap = make(map[string]interface{})
	json.Unmarshal(text, &base.ConfigMap)
}

//初始化全局数据、变量
func ini() {
	base.Con = base.NewConnect()
	base.CRAMs = base.NewCataRAMs()
	base.CRAMs.LoadCataRAM() //加载目录入内存

	base.Pcontent = base.Newcontent()
	base.Pcata = base.Newcata()
	base.PArticle = base.NewArticle()
	base.PSe = base.NewSe()

	pubgo.Tj = pubgo.Newtongji() //统计
}

//添加路由
func addrouters() {
	http.HandleFunc("/admin/", routers.AdminHtml)
	http.HandleFunc("/admin/cata/", routers.Catahtml)
	http.HandleFunc("/admin/search/", routers.Searchhtml)
	http.HandleFunc("/admin/art/", routers.Arthtml)

	http.HandleFunc("/api/cata/", routers.Cata)     //目录，get,post,put,delete
	http.HandleFunc("/api/art/", routers.Art)       //文章，get,post,put,delete
	http.HandleFunc("/api/search/", routers.Search) //搜索

	http.HandleFunc("/api/artitem/", routers.Artitem)       //获取目录下的文章列表
	http.HandleFunc("/api/idxfindpfx/", routers.Idxfindpfx) //搜索词为前缀的相关词
}

//运行服务
func run() {
	port := pubgo.ConfigMap["port"].(string) //从配置文件获取port
	err := http.ListenAndServe(":"+port, nil)
	//log.Fatal(err)
	if err != nil {
		fmt.Println("请更正错误后重启程序：", err)
	}
}
func main() {
	readconfig()
	ini()
	addrouters()
	//Test() //测试代码
	run()
}
