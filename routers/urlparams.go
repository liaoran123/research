package routers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// 用户提交的数据参数转map
func splitdata(reqdata string) map[string]string {
	params := make(map[string]string)
	pds := strings.Split(reqdata, "&")
	var ap []string
	var pitem, val string
	for _, p := range pds {
		ap = strings.Split(p, "=")
		//url.QueryUnescape(string(con))
		if len(ap) < 2 {
			continue
		}
		pitem, _ = url.QueryUnescape(ap[0])
		val, _ = url.QueryUnescape(ap[1])
		params[pitem] = val

	}
	return params
}

// 获取post的数据
func getPostData(req *http.Request) string {
	defer req.Body.Close()
	con, _ := ioutil.ReadAll(req.Body)
	//pdata, _ := url.QueryUnescape(string(con))
	//return pdata
	return string(con)
}

// 获取get的数据
func getGettData(req *http.Request) string {
	if strings.Contains(req.RequestURI, "?") {
		return strings.Split(req.RequestURI, "?")[1]
	}
	return ""
	//pdata, _ := url.QueryUnescape(string(pas))
	//return pdata
}

// 将post数据转为map
func postparas(req *http.Request) map[string]string {
	return splitdata(getPostData(req))
}

// 将get数据转为map
func getparas(req *http.Request) map[string]string {
	return splitdata(getGettData(req))
}
