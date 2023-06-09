package pubgo

//公共函数库
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ConfigMap map[string]interface{} //配置文件
)

// --截取前 l 个长度字符串
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

// 拼接多个[]byte
func JoinBytes(pBytes ...[]byte) []byte {
	len := len(pBytes)
	var buffer bytes.Buffer
	for index := 0; index < len; index++ {
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}

// 将关键词按非中文分解成数组
func GetKeys(word string) []string {
	reg := regexp.MustCompile(`[\p{Han}]+`) // 查找连续的汉字
	kws := reg.FindAllString(word, -1)      //,并生成数组
	return kws
}

// 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

// 获取程序绝对路径目录
func GetCurrentAbPath() string {
	exePath, err := os.Executable()
	if err != nil {
		//log.Fatal(err)
		return ""
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res + "\\"

}
func Reverse(str string) string {
	text := []rune(str)
	length := len(text)
	var result []rune
	for i := 0; i < length; i++ {
		result = append(result, text[length-i-1])
	}
	return string(result)
}

/*
//Float64ToByte Float64转byte

	func Float64ToByte(float float64) []byte {
		bits := math.Float64bits(float)
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, bits)
		return bytes
	}

//ByteToFloat64 byte转Float64

	func ByteToFloat64(bytes []byte) float64 {
		bits := binary.LittleEndian.Uint64(bytes)
		return math.Float64frombits(bits)
	}

////int转byte

	func IntToBytes(n int) []byte {
		data := int64(n)
		bytebuf := bytes.NewBuffer([]byte{})
		binary.Write(bytebuf, binary.BigEndian, data)
		return bytebuf.Bytes()
	}

//byte转int

	func BytesToInt(bys []byte) int {
		bytebuff := bytes.NewBuffer(bys)
		var data int64
		binary.Read(bytebuff, binary.BigEndian, &data)
		return int(data)
	}
*/
func Chekerr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
