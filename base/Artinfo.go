package base

//文章信息
type Artinfo struct {
	Title string `json:"title"`
	Split string `json:"split"`
	Url   string `json:"url"`
	Fid   int    `json:"fid"`
	Text  string `json:"text"`
}
