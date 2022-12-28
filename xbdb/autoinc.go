package xbdb

import (
	"strconv"
)

//自动增值
type Autoinc struct {
	id int
}

func NewAutoinc(id int) *Autoinc {
	return &Autoinc{
		id: id,
	}
}

//获取id并同时增值1
func (a *Autoinc) Getid() int {
	//a.mu.Lock() //由于leveldb是单线程，故不在这里加锁
	//defer a.mu.Unlock()
	id := a.id
	a.id++
	return id
}
func (a *Autoinc) GetidStr() string {
	return strconv.Itoa(a.Getid())
}

//获取当前id
func (a *Autoinc) GetidNoInc() int {
	id := a.id
	return id
}

func (a *Autoinc) GetidStrNoInc() string {
	return strconv.Itoa(a.GetidNoInc())
}
