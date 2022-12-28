package routers

import (
	"bytes"
	"research/xbdb"
)

//文章id对应的fid即目录id
var Artfid map[uint32]uint32

func LoadartRAM() {
	key := []byte("art" + xbdb.Split)
	tbd := Table["art"].Select.FindPrefix(key, true, 0, -1)
	if tbd != nil {
		toRAM(tbd)
	}
}
func toRAM(tbd *xbdb.TbData) {
	Artfid = make(map[uint32]uint32)
	var key [][]byte
	var artid, fid int
	for _, v := range tbd.Rd {
		key = bytes.Split(v, []byte(xbdb.Split))
		artid = xbdb.BytesToInt(key[0])
		fid = xbdb.BytesToInt(key[2])
		Artfid[uint32(artid)] = uint32(fid)
	}
}
