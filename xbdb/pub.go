// 小白数据库
// 表信息
package xbdb

import (
	"bytes"
	"encoding/binary"
	"math"
)

var (
	//xb  *leveldb.DB
	err error
)

type ReInfo struct {
	Succ   bool   `json:"Succ"` //是否成功
	Info   string `json:"Info"` //返回信息
	LastId string `json:"LastId"`
	//添加记录的id，主要是记录自动增值的id。由于有非自动增值的情况，故而使用string，可以兼容。
	Count int `json:"Count"` //修改、删除的记录数
}

// 拼接多个[]byte
func JoinBytes(pBytes ...[]byte) (r []byte) {
	r = bytes.Join(pBytes, []byte(""))
	return
}

// 拼接多个[]byte，遇到Nil则不再拼接
func JoinBytesNoNil(pBytes ...[]byte) []byte {
	blen := len(pBytes)
	var buffer bytes.Buffer
	for index := 0; index < blen; index++ {
		if len(pBytes[index]) == 0 {
			break
		}
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}

// 整型转换成字节
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
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Float64ToByte Float64转byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

// ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// Float64ToByte Float64转byte
func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

// ByteToFloat64 byte转Float64
func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

// 判断数组元素是否存在
func arryexist(a []string, v string) (r bool) {
	for _, av := range a {
		if v == av {
			r = true
			return
		}
	}
	return
}

// 将包括分隔符的数据转义
func SplitToCh(k []byte) (r []byte) {
	r = bytes.Replace(k, []byte(Split), []byte(ChSplit), -1)
	r = bytes.Replace(r, []byte(IdxSplit), []byte(ChIdxSplit), -1)
	return
}

// 将记录分解并将转义数据恢复
// 由于int，float转byte，会占用所有的特殊字符
// 添加时的“转义”只是标识，SplitRd根据标识转换，才能得到正确的结果。
// 所有数据都要经过此函数，方能得到未转义前的正确数据。
func SplitRd(Rd []byte) (r [][]byte) {
	csp := "[fgf0]"
	csp1 := "[fgf1]"
	rds := bytes.Replace(Rd, []byte(ChSplit), []byte(csp), -1)
	rds = bytes.Replace(rds, []byte(ChIdxSplit), []byte(csp1), -1)
	r = bytes.Split(rds, []byte(Split))
	for i, v := range r {
		r[i] = bytes.Replace(v, []byte(csp), []byte(Split), -1)
		r[i] = bytes.Replace(r[i], []byte(csp1), []byte(IdxSplit), -1)
	}
	return
}

/*
// 将包括分隔符的转义数据恢复
func ChToSplit(k []byte) (r []byte) {
	r = bytes.Replace(k, []byte(ChSplit), []byte(Split), -1)
	r = bytes.Replace(r, []byte(ChIdxSplit), []byte(IdxSplit), -1)
	return
}
*/
