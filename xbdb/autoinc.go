package xbdb

import (
	"strconv"
	"sync"
)

// 自动增值
type Autoinc struct {
	id int
	mu sync.RWMutex
}

func NewAutoinc(id int) *Autoinc {
	return &Autoinc{
		id: id,
	}
}

// 获取id并同时增值1
func (a *Autoinc) Getid() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	id := a.id
	a.id++
	return id
}
func (a *Autoinc) GetidStr() string {
	return strconv.Itoa(a.Getid())
}

// 获取当前id
func (a *Autoinc) GetidNoInc() int {
	id := a.id
	return id
}

// 获取未增值前的id值
func (a *Autoinc) GetidDic() int {
	id := a.id - 1
	return id
}
func (a *Autoinc) GetidStrNoInc() string {
	return strconv.Itoa(a.GetidNoInc())
}
