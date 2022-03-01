package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// PKCS5Padding 填充模式
func PKCS5Padding(data []byte, blockSize int) []byte {
	if blockSize != 8 {
		panic("wrong blocksize!")
	}
	padding := blockSize - len(data)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS5UnPadding 填充的反向操作,删除填充字符串
func PKCS5UnPadding(data []byte) []byte {
	//获取数据长度
	length := len(data)
	if length == 0 {
		panic("wrong data!")
	} else {
		//获取填充字符串长度
		unpadding := int(data[length-1])
		//截取切片,删除填充字节,并且返回明文
		return data[:(length - unpadding)]
	}
}

// PKCS7Padding 填充模式
func PKCS7Padding(data []byte, blockSize int) []byte {
	if blockSize < 0 || blockSize > 255 {
		panic("wrong blocksize!")
	}
	padding := blockSize - len(data)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7UnPadding 填充的反向操作,删除填充字符串
func PKCS7UnPadding(data []byte) []byte {
	//获取数据长度
	length := len(data)
	if length == 0 {
		panic("wrong data!")
	} else {
		//获取填充字符串长度
		unpadding := int(data[length-1])
		//截取切片,删除填充字节,并且返回明文
		return data[:(length - unpadding)]
	}
}

// ZeroPadding 等效于PKCS5
func ZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(data, padtext...)
}

// ZeroUnPadding 等效于PKCS5
func ZeroUnPadding(data []byte) []byte {
	return bytes.TrimFunc(data,
		func(r rune) bool {
			return r == rune(0)
		})
}

// AesEncrypt ...
func AesEncrypt(data string, key string) (string, error) {
	runeKey := []rune(key)
	// strings.Count(key, "") - 1
	if len(runeKey) != 16 && len(runeKey) != 24 && len(runeKey) != 32 {
		return "", nil
	}
	dataByte := []byte(data)
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	dataByte = PKCS7Padding(dataByte, blockSize)
	crypted := make([]byte, len(dataByte))
	blockMode.CryptBlocks(crypted, dataByte)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

// AesDecrypt ...
func AesDecrypt(crypted string, key string) (string, error) {
	runeKey := []rune(key)
	if len(runeKey) != 16 && len(runeKey) != 24 && len(runeKey) != 32 {
		return "", nil
	}
	cryptdByte, _ := base64.StdEncoding.DecodeString(crypted)
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, keyByte[:blockSize])
	dataByte := make([]byte, len(cryptdByte))
	blockMode.CryptBlocks(dataByte, cryptdByte)
	dataByte = PKCS7UnPadding(dataByte)
	if err != nil {
		return "", err
	}
	return string(dataByte), nil
}

// error
// // Base64Encrypt 16,24,32位字符串key，分别对应AES-128，AES-192，AES-256 加密方法
// func Base64Encrypt(data []byte, key []byte) ([]byte, error) {
// 	coder := base64.NewEncoding(string(key))
// 	return []byte(coder.EncodeToString(data)), nil
// 	// result, err := AesEncrypt(data, key)
// 	// if err != nil {
// 	// 	return nil, nil
// 	// }
// 	// return []byte(base64.StdEncoding.EncodeToString(result)), nil
// }

// // Base64Decrypt 16,24,32位字符串key，分别对应AES-128，AES-192，AES-256 加密方法
// func Base64Decrypt(crypted []byte, key []byte) ([]byte, error) {
// 	coder := base64.NewEncoding(string(key))
// 	return coder.DecodeString(string(crypted))
// 	// cryptdByte, err := base64.StdEncoding.DecodeString(string(crypted))
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// dataByte, err := AesDecrypt(cryptdByte, key)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// return dataByte, nil
// }

// Base64Encrypt ...
func Base64Encrypt(data []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(data)), nil
}

// Base64Decrypt ...
func Base64Decrypt(crypted []byte) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

// DesEncrypt 采用PKCS5
func DesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// pdData := ZeroPadding(data, blockSize)
	pdData := PKCS5Padding(data, blockSize)
	if len(pdData)%blockSize != 0 {
		return nil, errors.New("Need a multiple of the blocksize")
	}
	result := make([]byte, len(pdData))
	dst := result
	for len(pdData) > 0 {
		block.Encrypt(dst, pdData[:blockSize])
		pdData = pdData[blockSize:]
		dst = dst[blockSize:]
	}
	return []byte(hex.EncodeToString(result)), nil
}

// DesDecrypt 采用PKCS5
func DesDecrypt(crypted []byte, key []byte) ([]byte, error) {
	data, err := hex.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	result := make([]byte, len(data))
	dst := result
	if len(data)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(data) > 0 {
		block.Decrypt(dst, data[:blockSize])
		data = data[blockSize:]
		dst = dst[blockSize:]
	}
	// result = ZeroUnPadding(result)
	result = PKCS5UnPadding(data)
	return result, nil
}

// DesCBCEncrypt ...
func DesCBCEncrypt(data []byte, key []byte, ivb []byte) ([]byte, error) {
	var iv []byte
	if ivb == nil {
		iv = key
	} else if len(key) != len(ivb) {
		return nil, errors.New("The length of iv must be the same as the Block's block!")
	} else {
		iv = ivb
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data = PKCS5Padding(data, block.BlockSize())
	crypted := make([]byte, len(data))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(crypted, data)
	return []byte(base64.StdEncoding.EncodeToString(crypted)), nil
}

// DesCBCDecrypt ...
func DesCBCDecrypt(crypted []byte, key []byte, ivb []byte) ([]byte, error) {
	var iv []byte
	if ivb == nil {
		iv = key
	} else if len(key) != len(ivb) {
		return nil, errors.New("The length of iv must be the same as the Block's block!")
	} else {
		iv = ivb
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cryptedBase64, err := base64.StdEncoding.DecodeString(string(crypted))
	if err != nil {
		return nil, err
	}
	data := make([]byte, len(cryptedBase64))
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(data, cryptedBase64)
	data = PKCS5UnPadding(data)
	return data, nil
}

// TripleDesEncrypt ...
func TripleDesEncrypt(data []byte, key []byte, ivb []byte) ([]byte, error) {
	var iv []byte
	if len(key) != 24 {
		return nil, errors.New("The length of key must be:24!")
	}
	if ivb == nil {
		iv = key[:8]
	} else if len(ivb) == 8 {
		iv = ivb
	} else {
		return nil, errors.New("The length of iv must be:8!")
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	data = PKCS7Padding(data, block.BlockSize())
	// blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	return crypted, nil
}

// TripleDesDecrypt ...
func TripleDesDecrypt(crypted []byte, key []byte, ivb []byte) ([]byte, error) {
	var iv []byte
	if len(key) != 24 {
		return nil, errors.New("The length of key must be:24!")
	}
	if ivb == nil {
		iv = key[:8]
	} else if len(ivb) == 8 {
		iv = ivb
	} else {
		return nil, errors.New("The length of iv must be:8!")
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	// blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	blockMode := cipher.NewCBCDecrypter(block, iv)
	dataByte := make([]byte, len(crypted))
	blockMode.CryptBlocks(dataByte, crypted)
	dataByte = PKCS7UnPadding(dataByte)
	return dataByte, nil
}

// RsaEncrypt ...
func RsaEncrypt(data []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// RsaDecrypt ...
func RsaDecrypt(crypted []byte, privateKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, crypted)
}

// MD5Sum ...
func MD5Encrypt(data string) string {
	h := md5.New()
	if _, err := h.Write([]byte(data)); err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

// MD5Check crypted-密文,data-明文
func MD5Check(crypted string, data string) bool {
	return strings.EqualFold(MD5Encrypt(data), crypted)
}

// SHA1EncryptHex ...
func SHA1EncryptHex(data string) string {
	h := sha1.New()
	if _, err := h.Write([]byte(data)); err != nil {
		return ""
	}
	dataByte := h.Sum(nil)
	return fmt.Sprintf("%x", dataByte)
}

func MD5Sum(file *os.File) string {
	h := md5.New()
	_, _ = io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}

func SHA1Sum(file *os.File) string {
	h := sha1.New()
	_, _ = io.Copy(h, file)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5SumByFilePath(filePath string) string {
	b, _ := ioutil.ReadFile(filePath)
	return fmt.Sprintf("%x", md5.Sum(b))
}

func SHA1SumByFilePath(filePath string) string {
	b, _ := ioutil.ReadFile(filePath)
	return fmt.Sprintf("%x", sha1.Sum(b))
}
