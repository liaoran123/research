package main

import (
	"fmt"
	"research/base"
	"research/pubgo"
	"strings"
)

//-------测试代码-----------
func Test() {
	fmt.Println("-------目录-----------")
	base.Con.Getartdb().FindPrefixFun("ca-", print)

	fmt.Println("-------目录fid索引-----------")
	base.Con.Getartdb().FindPrefixFun("cf-", print1)

	fmt.Println("-------目录fid对应文章id-----------")
	base.Con.Getartdb().FindPrefixFun("fa-", print1)

	fmt.Println("-------文章id对应目录fid-----------")
	base.Con.Getartdb().FindPrefixFun("af-", print1)

	fmt.Println("-------文章-----------")
	base.Con.Getartdb().FindPrefixFun("c-", print)

	fmt.Println("-------文章索引-----------")
	base.Con.Getartdb().FindPrefixFun("i-", print)

}
func print(k, v []byte) {
	fmt.Println(string(k), ",", string(v))
}
func print1(k, v []byte) {
	ks := strings.Split(string(k), "-")
	k1 := pubgo.BytesToInt([]byte(ks[1]))
	k2 := pubgo.BytesToInt([]byte(ks[2]))
	fmt.Println(ks, k1, k2)
}
