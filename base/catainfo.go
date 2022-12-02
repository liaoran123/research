package base

//目录信息
type CataInfo struct {
	Id     int    `json:"id"`     //目录id
	Fid    int    `json:"fid"`    //目录fid
	Name   string `json:"title"`  //目录名称
	Isleaf string `json:"Isleaf"` //是否叶子目录

}

//所有子目录信息
type CataInfos struct {
	Catas []CataInfo `json:"Catainfo"`
}
