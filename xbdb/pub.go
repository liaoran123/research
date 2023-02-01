//小白数据库
//表信息
package xbdb

import (
	"bytes"
	"encoding/binary"
	"math"
)

var (
	err error
)

type ReInfo struct {
	Succ   bool   `json:"Succ"`
	Info   string `json:"Info"`
	LastId string `json:"LastId"`
	//添加记录的id，主要是记录自动增值的id。由于有非自动增值的情况，故而使用string，可以兼容。
	Count int `json:"Count"` //修改、删除的记录数
}

//拼接多个[]byte
func JoinBytes(pBytes ...[]byte) (r []byte) {
	r = bytes.Join(pBytes, []byte(""))
	return
}

/*
//拼接多个[]byte
func JoinBytes(pBytes ...[]byte) []byte {
	len := len(pBytes)
	var buffer bytes.Buffer
	for index := 0; index < len; index++ {
		buffer.Write(pBytes[index])
	}
	return buffer.Bytes()
}
*/

//整型转换成字节
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
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

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

//Float64ToByte Float64转byte
func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

//ByteToFloat64 byte转Float64
func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}
