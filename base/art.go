package base

import (
	"strconv"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb/util"
)

//搜索系统文章的添加、修改、删除
var PArticle *Article

type Article struct {
	cont *content
	//idx  *idx
	cata *cata
	//fididx *fididx
	//sp []string
	mu sync.RWMutex
}

func NewArticle() *Article {
	return &Article{
		cont: Pcontent, //Newcontent(),
		//idx:  Newidx(),
		cata: Pcata, //Newcata(),
		//fididx: Newfididx(),
	}
}

//添加。支持手动回滚。
func (a *Article) Insert(id int, title, text, split, url string, fcataid int) (r bool) {
	a.mu.Lock() //实际上leveldb只支持单进程。
	defer a.mu.Unlock()
	r = true
	r = r && a.cont.Insert(id, fcataid, title, text, split, url) //添加内容
	//r = r && a.cont.ainc.Writelastid()                   //事务出问题，不会写入自动增值表的最后id。单进程下可以支持回滚。用户重新写入，数据即可保证完整性。
	return
}

//删除
func (a *Article) Delete(id int) (r bool) {
	a.mu.Lock() //实际上leveldb只支持单进程。
	defer a.mu.Unlock()
	r = true
	r = r && a.cont.Delete(id)
	//r = r && a.cont.ainc.Writelastid()                   //事务出问题，不会写入自动增值表的最后id。单进程下可以支持回滚。用户重新写入，数据即可保证完整性。
	return
}

//修改=删除+添加
func (a *Article) Put(id int, title, text, split, url string, fcataid int) (r bool) {
	r = a.Delete(id)
	r = r && a.cont.InsertId(id, fcataid, title, text, split, url)
	return
}

//-获取某目录下的文章列表
func (a *Article) GetCataArt(fcataid, count int, p string) (r RArtItem) {
	ct := count
	if ct > 108 || ct == 0 { //每页最大和默认108条记录
		ct = 108
	}
	Prefix := JoinBytes([]byte(a.cont.fidx.tbn+"~"), IntToBytes(fcataid), []byte("~"))
	iter := Con.Getartdb().Db.NewIterator(util.BytesPrefix([]byte(Prefix)), nil)
	var b bool
	if p == "" {
		b = iter.First()
	} else { //重新组成key
		ps := strings.Split(p, "~")
		cid, _ := strconv.Atoi(ps[0])
		aid, _ := strconv.Atoi(ps[1])
		pkey := JoinBytes([]byte(a.cont.fidx.tbn+"~"), IntToBytes(cid), []byte("~"), IntToBytes(aid))
		b = iter.Seek([]byte(pkey))
	}
	var artid, cataid, loop int
	var artkey, sec []byte
	var lastkey string
	var ks []string
	var ai []ArtItem
	artItem := ArtItem{}
	for b {
		ks = strings.Split(string(iter.Key()), "~")
		artid = BytesToInt([]byte(ks[2]))
		cataid = BytesToInt([]byte(ks[1]))
		artkey = JoinBytes([]byte(a.cont.tbn+"~"), IntToBytes(artid), []byte("~"), IntToBytes(0)) //文章的第0个段落就文章标题
		sec, _ = Con.Getartdb().Db.Get(artkey, nil)
		/*
			if string(sec) == "" { //被删除的不显示
				b = iter.Next()
				continue
			}*/
		artItem.Id = artid
		artItem.Title = string(sec)
		ai = append(ai, artItem)
		lastkey = strconv.Itoa(cataid) + "~" + strconv.Itoa(artid)
		loop++
		if loop >= ct { //每页最多返回108条记录。必须限制，以免数据量大拖垮资源。
			break
		}
		b = iter.Next()
	}
	Release(iter)
	r.ArtItems = ai
	r.LastKey = lastkey
	return
}

type ArtItem struct {
	Id    int    `json:"id"`    //文章id
	Title string `json:"title"` //文章标题
}
type RArtItem struct {
	ArtItems []ArtItem `json:"ArtItem"`
	LastKey  string    `json:"p"`
}
