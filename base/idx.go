package base

/*
i- 内容索引表
.k=i-长度7的分段内容遍历分词-文章id-分段id;v=空值
.i-最后后面加内容表的唯一id："文章id-分段id"，为了相同分词之后按内容表先后排序一致
*/
import (
	"strings"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//索引表
type idx struct {
	tbn        string
	keywordlen int
	//disparte   []string
}

func Newidx() *idx {
	keywordlen := 7
	if ConfigMap["keywordlen"] != nil {
		keywordlen = int(ConfigMap["keywordlen"].(float64))
	}
	return &idx{
		tbn:        "i",
		keywordlen: keywordlen,
	}
}

//******添加**********

//添加内容索引.pfx=表名
//Insert
func (i *idx) Act(artid, secid, fcataid int, sec string, f func(kw string, k, v []byte) error) (r bool) {
	var pk, pv []byte
	disparte := i.ForDisparte(sec)
	r = true
	for _, dv := range disparte {
		if strings.TrimSpace(dv) == "" {
			continue
		}
		pk = JoinBytes([]byte(i.tbn+"-"), []byte(dv), []byte("-"), IntToBytes(artid), []byte("-"), IntToBytes(secid))
		pv = []byte{}
		if fcataid != 0 {
			/*
				//所属目录id存入value中，以通过fcataid进行目录区间搜索的判断。
				通过GetArtFCataId也可以获取fcataid，但是速度慢。
				通过直接存在pv中，则速度快。用空间换时间方式。
			*/
			pv = IntToBytes(fcataid)
		}
		if r {
			//err = Con.Getidxdb(dv).Db.Put(pk, pv, nil)
			err = f(dv, pk, pv)
		}
		r = r && err == nil
		//fmt.Println(string(pk), cataid, secid)
		Chekerr()
	}
	return
}
func (i *idx) Insert(kw string, k, v []byte) error {
	return Con.Getidxdb(kw).Db.Put(k, v, nil)
}
func (i *idx) Delete(kw string, k, v []byte) error {
	return Con.Getidxdb(kw).Db.Delete(k, nil)
}

//遍历分词
func (i *idx) ForDisparte(nr string) (disparte []string) {
	var knr string //, fid
	var ml, cl int
	var r, idxstr []rune
	r = []rune(nr)
	cl = len([]rune(nr))
	for cl > 0 {
		if cl >= i.keywordlen {
			ml = i.keywordlen
		} else {
			ml = cl
		}
		idxstr = r[:ml]
		knr = string(idxstr)
		disparte = append(disparte, knr)
		r = r[1:]
		cl = len(r)
	}
	return
}

//********查询********************
//获取pfx为前缀的查询索引游标。
func (i *idx) Getiter(pfx string) iterator.Iterator {
	return Con.Getidxdb(pfx).Db.NewIterator(util.BytesPrefix([]byte(i.tbn+"-"+pfx)), nil)
}

func (i *idx) GetPfx(pfx string, top int) (r []string) {
	//Con.Getidxdb(pfx).FindPrefixTopFun(pfx, top)
	if pfx == "" {
		return
	}
	max := top
	if max > 21 || max == 0 {
		max = 21 //最大默认是21
	}
	iter := Con.Getidxdb(pfx).Db.NewIterator(util.BytesPrefix([]byte(i.tbn+"-"+pfx)), nil)
	loop := 0
	var ks []string
	var keys string
	for iter.Next() {
		ks = strings.Split(string(iter.Key()), "-")
		if !strings.Contains(keys, ks[1]+"|") {
			keys += ks[1] + "|"
		} else {
			continue
		}
		loop++
		if loop > max {
			break
		}
	}
	Release(iter)
	r = strings.Split(keys, "|")
	return
}
