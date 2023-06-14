package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"research/pubgo"
	"research/routers"
)

// 读取配置文件装入map
func readconfig() {
	path := pubgo.GetCurrentAbPath()
	text, _ := ioutil.ReadFile(path + "routers/config.json")
	routers.ConfigMap = make(map[string]interface{})
	json.Unmarshal(text, &routers.ConfigMap)
}

// 初始化全局数据、变量
func ini() {
	routers.Ini()
	pubgo.Tj = pubgo.Newtongji() //统计
}

// 添加路由
func addrouters() {
	http.HandleFunc("/admin/", routers.AdminHtml)
	http.HandleFunc("/admin/cata/", routers.Catahtml)
	http.HandleFunc("/admin/search/", routers.Searchhtml)
	http.HandleFunc("/admin/art/", routers.Arthtml)

	http.HandleFunc("/api/cata/", routers.Cata) //目录，get,post,put,delete
	http.HandleFunc("/api/art/", routers.Art)   //文章，get,post,put,delete
	//对应blazor的GetFromJsonAsync的json格式
	http.HandleFunc("/api/search/", routers.Search) //http.HandleFunc("/api/search/", routers.Search) //搜索

	http.HandleFunc("/api/art/item/", routers.Artitem) //获取目录下的文章列表
	http.HandleFunc("/api/art/meta/", routers.Meta)    //获取文章摘录
	http.HandleFunc("/api/Idxpfx/", routers.Idxpfx)    //搜索词为前缀的相关词

	http.HandleFunc("/test", routers.Test)
}

// 运行服务
func run() {
	fmt.Println("ReSearch服务器程序启动!")
	port := routers.ConfigMap["port"].(string) //从配置文件获取port
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
