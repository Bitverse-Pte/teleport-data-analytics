package tools

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"time"
)

var randomStrSource = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandomStr(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	for i := 0; i < length; i++ {
		result[i] = randomStrSource[r.Intn(len(randomStrSource))]
	}
	return string(result)
}

func GetRandomNumStr(length int) string {
	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	for i := 0; i < length; i++ {
		result[i] = byte('0' + r.Intn(10)) //0 - 9
	}
	return string(result)
}

func Uint8ToBool(data uint8) bool {
	if data > 0 {
		return true
	}
	return false
}

func ToString(i interface{}) string {
	if b, err := json.Marshal(i); err == nil {
		return string(b)
	}
	return "failed..."
}

func MakeUidWithTime(prefix string, random_len int) string {
	const time_base_format = "060102030405"
	uid := prefix + time.Now().Format(time_base_format) + GetRandomNumStr(random_len)
	return uid
}

func GetSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
