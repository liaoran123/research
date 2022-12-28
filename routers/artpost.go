package routers

import (
	"encoding/json"
	"net/http"
	"research/xbdb"
	"strconv"
	"strings"
)

func artpost(w http.ResponseWriter, req *http.Request) {
	mu.Lock() //leveldb仅支持单进程数据操作。
	defer mu.Unlock()
	params := postparas(req)
	r := InsArt(params)
	json.NewEncoder(w).Encode(r)

}
func InsArt(params map[string]string) (r xbdb.ReInfo) {
	serpsw := ConfigMap["pws"].(string)                             //服务器端不设置密码，即不可以进行操作
	if serpsw == "" || params["psw"] != ConfigMap["pws"].(string) { //密码不对
		r.Info = "密码不正确！"
		return
	}
	//url参数名称必须与表字段名称一致
	//修改也是添加，这样可以得到搜索容错词.
	//需要同义词等搜索也可以使用这样的方式添加。
	//需要真正修改，则先手动删除再添加
	//频繁的修改，将会产生数据冗余
	r = InsOrUpd("art", params, "ins") //r := InsOrUpd("art", params, params["iou"])
	if !r.Succ {
		return
	}
	//上面是添加文章信息
	//下面是添加一对多的文章分解为多个段落
	id := params["id"]
	if id == "" {
		iid := Table["art"].Ac.GetidNoInc() - 1
		id = strconv.Itoa(iid)
	}
	//c表的id=文章id+句子段落id+上级目录id
	paramsc := map[string]string{
		"id": "",
		"s":  "",
	}
	//《"+title+"》\n 将标题加入内容，即不需要多加一个标题的搜索。通过“《”+关键词就能专门搜索标题
	text := "《" + params["title"] + "》\n" + params["text"]
	sectexts := TextSplit(text, params["split"])
	iv := ""
	i := 0
	var art int
	for _, v := range sectexts {
		iv = strings.TrimSpace(v)
		if iv == "" {
			continue
		}
		//文章id+句子段落id
		art, _ = strconv.Atoi(id)
		paramsc["id"] = ArtSecToId(art, i)
		paramsc["s"] = v
		//paramsc["pos"] = strconv.Itoa(j)
		r = InsOrUpd("c", paramsc, "ins")
		if !r.Succ {
			return
		}
		i++
	}
	return
}

//***********添加*******************
//添加内容，将文章分成多个句子段落后添加到表
func TextSplit(text, split string) (section []string) {
	itext := text //title+"~"+text
	if split != "" {
		//支持多个分段匹配标签。中文常见是“。”.
		//空格是组合查询，由于支持英文，故而空格不作默认分隔符
		ss := strings.Split(split, "|")
		for _, v := range ss {
			itext = strings.Replace(itext, v, v+"\n", -1) //分配段落
		}
	}
	section = strings.Split(itext, "\n")
	return
}

const idssplit = "||"

//artid+secid组合成一个字符串
func ArtSecToId(Art, Sec int) (r string) {
	as := xbdb.IntToBytes(Art)
	ss := xbdb.IntToBytes(Sec)
	r = string(as) + idssplit + string(ss) //数字转byte不会包含有A或其他字母
	return
}

//一个字符串分解为artid+secid
func IdToArtSec(r string) (Art, Sec int) {
	rs := strings.Split(r, idssplit)
	if len(rs) != 2 {
		return
	}
	Art = xbdb.BytesToInt([]byte(rs[0]))
	Sec = xbdb.BytesToInt([]byte(rs[1]))
	return
}
