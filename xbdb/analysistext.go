package xbdb

import (
	"strings"
	"unicode"
)

const noPattern = " "

//解析搜索词的匹配模型
/*
判断是否为中文：unicode.Is(unicode.Han, v) //unicode.Is(unicode.Scripts["Han"], v)  =1

判断是否为字母： unicode.IsLetter(v) =2

判断是否为数字： unicode.IsNumber(v) =3

判断是否为Unicode标点字符 :unicode.IsPunct(v) =4

判断是否为空白符号： unicode.IsSpace(v)

判断是否为十进制数字： unicode.IsDigir(v)
*/
//搜索词的匹配模型可以组合匹配。（patterns 可以是以上1--4的1个或多个组合。比如，中文和字母=[]int{1，2}）
func Analysis(kw string, patterns []int) []string {
	rkw, k := "", ""
	for _, v := range kw {
		for _, pt := range patterns {
			if pt == 0 { //0，表示全部原样返回
				rkw += string(v)
				break
			}
			k = IsRet(v, pt)
			if k != noPattern {
				rkw += k
				break
			}
			rkw += noPattern
		}
	}
	return strings.Split(rkw, " ")
}

// 是中文或字母等等即原样返回，否则返回空格
func IsRet(c rune, pattern int) string {
	if pattern == 1 {
		if unicode.Is(unicode.Han, c) { //返回中文
			return string(c)
		}
	}
	if pattern == 2 {
		if unicode.IsLetter(c) { //返回字母
			return string(c)
		}
	}
	if pattern == 3 {
		if unicode.IsNumber(c) { //返回数字
			return string(c)
		}
	}
	if pattern == 4 {
		if unicode.IsPunct(c) { //返回标点符号
			return string(c)
		}
	}
	if pattern == 5 {
		if isletter(c) { //返回自定义
			return string(c)
		}
	}
	return noPattern
}

// 自定义字符
func isletter(r rune) bool {
	lt := "《》" //自定义字符
	for _, v := range lt {
		if v == r {
			return true
		}
	}
	return false
}
