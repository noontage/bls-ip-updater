package randstr

import (
	"crypto/rand"
)

const letterTable = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func NewString(n int) string {
	return readRandomString(n)
}

func readRandomString(n int) string {
	buf := readRandom(n)
	for i, v := range buf {
		buf[i] = letterTable[v%byte(len(letterTable))]
	}
	return string(buf)
}

func readRandom(n int) []byte {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}
