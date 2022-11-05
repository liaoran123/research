package base

/*
var IdsPool = sync.Pool{
	New: func() interface{} {
		return new(Ids)
	},
}
*/
//返回结果集。文章和段落的id
type asid struct {
	cataid int
	artid  int
	secid  int
}

/*
func newasid(artid, secid int) asid {
	return asid{artid: artid, secid: secid}
}
*/
