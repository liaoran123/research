package base

//返回结果集。文章和段落的id
type Rset struct {
	//Cataid int `json:"Cataid"` //目录id
	//CataName  string `json:"CataName"`  //目录名称
	CataDir   []CataInfo `json:"CataDir"`   //目录路径
	Artid     int        `json:"Artid"`     //文章id
	Title     string     `json:"Title"`     //文章标题
	ArtUrl    string     `json:"ArtUrl"`    //文章网址
	Secid     int        `json:"Secid"`     //句子或段落开始id
	LastSecid int        `json:"LastSecid"` //句子或段落结束id
	Text      string     `json:"Text"`      //搜索结果节录。从Secid--LastSecid的内容相加。
}

/*type CataDir struct {
	Cataid   int    `json:"id"`    //目录id
	CataName string `json:"title"` //目录名称
}
*/
/*
func newRset(artid, secid int, text string) *Rset {
	return &Rset{
		Artid: artid,
		Secid: secid,
		Text:  text,
	}
}
*/

type RsetAll struct {
	Set     []Rset `json:"Rset"`    //文章id
	Lastkey string `json:"p"`       //该页最后一个key值，作为下一页的起始位置。
	SeTime  string `json:"SeTime"`  //搜索用时
	SetTime string `json:"SetTime"` //数据集结用时
}

func (R *RsetAll) Reset() {
	R.Set = R.Set[:0]
	R.Lastkey = ""
	R.SeTime = ""
	R.SetTime = ""
}
