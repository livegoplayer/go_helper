package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type HashFunc func(str string) int64

var (
	BKDRHashFunc HashFunc = BKDRHash
	APHashFunc   HashFunc = APHash
	SDBMHashFunc HashFunc = SDBMHash
	RSHashFunc   HashFunc = RSHash
	JSHashFunc   HashFunc = JSHash
	ELFHashFunc  HashFunc = ELFHash
	DJBHashFunc  HashFunc = DJBHash
)

func BKDRHash(str string) int64 {
	const seed int64 = 131
	r := []rune(str)
	count, size := len(r), len(r)
	var hash int64
	for count > 0 {
		hash = hash*seed + int64(r[size-count])
		count--
	}
	return hash & 0x7FFFFFFF
}

func APHash(str string) int64 {
	r := []rune(str)
	count := len(r)
	var hash int64
	for i := 0; i < count; i++ {
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ int64(r[i]) ^ (hash >> 3)
		} else {
			hash ^= ^((hash << 11) ^ int64(r[i]) ^ (hash >> 5))
		}
		count--
	}
	return hash & 0x7FFFFFFF
}

func SDBMHash(str string) int64 {
	r := []rune(str)
	count, size := len(r), len(r)
	var hash int64
	for count > 0 {
		hash = int64(r[size-count]) + (hash << 6) + (hash << 16) - hash
		count--
	}
	return hash & 0x7FFFFFFF
}

func RSHash(str string) int64 {
	var a int64 = 63689
	const b int64 = 378551
	r := []rune(str)
	count, size := len(r), len(r)
	var hash int64
	for count > 0 {
		hash = hash*a + int64(r[size-count])
		a *= b
		count--
	}
	return hash & 0x7FFFFFFF
}

func JSHash(str string) int64 {
	var hash int64 = 1315423911
	r := []rune(str)
	count, size := len(r), len(r)
	for count > 0 {
		hash ^= (hash << 5) + int64(r[size-count]) + (hash >> 2)
		count--
	}
	return hash & 0x7FFFFFFF
}

func ELFHash(str string) int64 {
	var x, hash int64
	r := []rune(str)
	count, size := len(r), len(r)
	for count > 0 {
		hash = (hash << 4) + int64(r[size-count])
		if x = hash & 0xF0000000; x != 0 {
			hash ^= x >> 24
			hash &= ^x
		}
		count--
	}
	return hash & 0x7FFFFFFF
}

func DJBHash(str string) int64 {
	r := []rune(str)
	count, size := len(r), len(r)
	var hash int64
	for count > 0 {
		hash += (hash << 5) + int64(r[size-count])
		count--
	}
	return hash & 0x7FFFFFFF
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Sha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - (len(ciphertext) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesCbcEncrypt(str string, key []byte, IV []byte) ([]byte, error) {
	origData := []byte(str)

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := cipherBlock.BlockSize()
	origData = pKCS5Padding(origData, blockSize)

	crypted := make([]byte, len(origData))
	cipher.NewCBCEncrypter(cipherBlock, IV).CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesCbcDecrypt(encrypted []byte, key []byte, IV []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypt := make([]byte, len(encrypted))

	cipher.NewCBCDecrypter(cipherBlock, IV).CryptBlocks(decrypt, encrypted)
	return decrypt, nil
}

func GetHashCode(str string, funcType ...HashFunc) int64 {
	hashCode := BKDRHashFunc
	if len(funcType) > 0 {
		hashCode = funcType[0]
	}
	return hashCode(str)
}
