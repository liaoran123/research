package base

//公共函数库
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var (
	ConfigMap map[string]interface{} //配置文件
)

//--截取前 l 个长度字符串
func Sublen(str string, L int) string {
	runek := []rune(str) //包含中文必须如此才能得到正确的长度
	k := ""
	if len(runek) > L {
		k = string(runek[:L]) //截取35位
	} else {
		k = str
	}
	return k
}

//拼接多个[]byte
func JoinBytes(pBytes ...[]byte) []byte {
	len := len(pBytes)
	var buffer bytes.Buffer
	for index := 0; index < len; index++ {
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}

//将关键词按非中文分解成数组
func GetKeys(word string) []string {
	reg := regexp.MustCompile(`[\p{Han}]+`) // 查找连续的汉字
	kws := reg.FindAllString(word, -1)      //,并生成数组
	return kws
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

func Chekerr() {
	if err != nil {
		fmt.Println(err)
	}
}

func Release(iter iterator.Iterator) {
	iter.Release()
	err = iter.Error()
	Chekerr()
}
