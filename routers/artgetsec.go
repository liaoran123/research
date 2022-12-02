package routers

import (
	"encoding/json"
	"net/http"
	"research/base"
	"strconv"
)

//获取文章一段落内容
func artgetsec(w http.ResponseWriter, req *http.Request) {
	params := getparas(req)
	id := params["id"]
	secid := params["secid"]
	iid, _ := strconv.Atoi(id)
	isecid, _ := strconv.Atoi(secid)
	r := base.Pcontent.GetOneSec(iid, isecid)
	st := sectext{}
	if r != nil {
		st.Text = string(r)
		//st.Text = strings.Replace(st.Text, "﹣", "-", -1) //text = strings.Replace(text, "-", "﹣", -1) //-是系统保留字，需要转义为﹣。
	} else {
		st.Text = "【已结束】"
	}

	json.NewEncoder(w).Encode(st)
}

type sectext struct {
	Text string `json:"text" `
}
