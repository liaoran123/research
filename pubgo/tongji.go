package pubgo

import "time"

//---统计每天流量
var Tj *tongji

type tongji struct {
	Tjs map[string]*daybws
}

func Newtongji() *tongji {
	return &tongji{map[string]*daybws{}}
}
func (t *tongji) Brows(dir string) {
	if dir == "" {
		dir = "home"
	}
	if db, ok := t.Tjs[dir]; ok {
		db.Add()
	} else {
		t.Tjs[dir] = Newdaybws()
	}
}

type daybws struct {
	day int
	Bws int //--每天流量
}

func Newdaybws() *daybws {
	return &daybws{Bws: 1}
}
func (d *daybws) Add() { //按天统计点击
	day := time.Now().Day()
	if day == d.day {
		d.Bws++
	} else {
		d.day = day
		d.Bws = 1
	}
}
