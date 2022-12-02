package base

//文章所属目录id，即fid索引表，一个fid对多个artid
type fididx struct {
	tbn string
	//rfi *Refididx
}

func Newfididx() *fididx {
	return &fididx{
		tbn: "fa",
		//rfi: NewRefididx(),
	}
}

//fid索引，通过fid找到artid
//添加
func (f *fididx) Insert(fid, artid int) (r bool) {
	//JoinBytes([]byte(f.tbn+"~"), IntToBytes(fid), []byte("~"), IntToBytes(artid))
	err = Con.Getartdb().Db.Put(f.setkey(fid, artid), []byte{}, nil)
	Chekerr()
	r = err == nil
	//r = r && f.rfi.Insert(artid, fid)
	return
}

//删除
func (f *fididx) Delete(fid, artid int) (r bool) {
	err = Con.Getartdb().Db.Delete(f.setkey(fid, artid), nil)
	Chekerr()
	r = err == nil
	//r = r && f.rfi.Delete(artid, fid)
	return
}

//key=fa-fid-artid
func (f *fididx) setkey(fid, artid int) (r []byte) {
	r = JoinBytes([]byte(f.tbn+"~"), IntToBytes(fid), []byte("~"), IntToBytes(artid))
	return
}
