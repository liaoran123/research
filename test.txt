package main

import (
	"fmt"
	"research/routers"
	"research/xbdb"
)

//-------测试代码-----------

func Test() {

	routers.Table["ca"].Select.ForDb(print)
	r := routers.Table["ca"].FindPrefix([]byte("ca"+xbdb.Split), true, 0, -1)
	fmt.Println("ca", r)
	r.Release()
	r = routers.Table["art"].FindPrefix([]byte("art"+xbdb.Split), true, 0, -1)
	fmt.Println("art", r)
	r.Release()
	r = routers.Table["c"].FindPrefix([]byte("c"+xbdb.Split), true, 0, -1)
	fmt.Println("c", r)
	r.Release()

}

func print(rd []byte) bool {
	fmt.Println(string(rd))
	return true
}

/*
func chcata(k, v []byte) {
	key := string(k)
	id := strings.Split(key, "~")[1]
	iid := base.BytesToInt([]byte(id))
	//fmt.Printf("iid: %v\n", iid)
	//fmt.Println(key, ",", string(v), strings.Trim(string(v), "0"), strings.Trim(string(v), "0")+"1", "---------------")

	//--没有子目录即是叶子目录
	Prefix := base.JoinBytes([]byte("cf~"), base.IntToBytes(iid), []byte("~"))
	iter := db.Db.NewIterator(util.BytesPrefix([]byte(Prefix)), nil)
	loop := 0
	for iter.Next() {
		//fmt.Println(string(iter.Key()), ",", string(iter.Value()))
		loop++
		break
	}
	if loop == 0 {
		//fmt.Println("isleaf")
		fmt.Println(key, ",", string(v))
		db.Db.Put(k, []byte(strings.Trim(string(v), "0")+"1"), nil)
	}
}


func print1(k, v []byte) {
	ks := strings.Split(string(k), "\n")
	k1 := pubgo.BytesToInt([]byte(ks[1]))
	k2 := pubgo.BytesToInt([]byte(ks[2]))
	fmt.Println(ks, k1, k2)
}
*/
