package utils

import (
	"bytes"
	"math"
)

type Enc62 struct {
	Alphabet string
	Base     int
}

func NewEnc62(alphabet string) *Enc62 {
	if alphabet == "" {
		alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	return &Enc62{Alphabet: alphabet, Base: len(alphabet)}
}

func (e *Enc62) Encode(num int64) string {
	length := int64(e.Base)
	result := ""
	for num > 0 {
		remainder := num % length
		result = string(e.Alphabet[remainder]) + result
		num /= length
	}
	return result
}

func (e *Enc62) Decode(str string) int64 {
	number := int64(0)
	idx := 0.0
	chars := []byte(e.Alphabet)

	charsLen := float64(e.Base)
	strLen := float64(len(str))
	for _, c := range []byte(str) {
		power := strLen - (idx + 1)
		index := int64(bytes.IndexByte(chars, c))
		number += index * int64(math.Pow(charsLen, power))
		idx++
	}
	return number
}
