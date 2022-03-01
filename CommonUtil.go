package util

import (
	"bytes"
	"encoding/binary"
	"os"
)

// IntToBytes 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt 字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func FileOrPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func MapContainStringString(m map[string]string, key string) bool {
	if _, exist := m[key]; exist {
		return true
	}
	return false
}

func MapContainInterfaceString(m map[string]interface{}, key string) bool {
	if _, exist := m[key]; exist {
		return true
	}
	return false
}

func SilceContainStringString(strs []string, name string) bool {
	for _, str := range strs {
		if name == str {
			return true
		}
	}
	return false
}

func ContentToLine(content []byte) []string {
	var lines []string
	var line []byte
	for _, b := range content {
		if b == '\n' {
			lines = append(lines, string(line))
			line = line[0:0]
			continue
		}
		line = append(line, b)
	}
	return lines
}
