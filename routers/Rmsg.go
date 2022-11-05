package routers

//添加删除修改返回信息
type Rmsg struct {
	Msg  string `json:"Msg"`
	Succ bool   `json:"Succ"`
	//Time string `json:"time"`
}

func NewRmsg() Rmsg {
	return Rmsg{}
}
