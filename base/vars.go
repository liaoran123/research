package base

import (
	"bytes"
	"sync"
)

var (
	//Art *Article
	err error
)

//sync.Pool类似公共变量的概念。
//.Get()加锁；.Put()解锁。
var RsetAllPool = sync.Pool{
	New: func() interface{} {
		return new(RsetAll)
	},
}
var RsetPool = sync.Pool{
	New: func() interface{} {
		return new(Rset)
	},
}
var ArtBuf = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
