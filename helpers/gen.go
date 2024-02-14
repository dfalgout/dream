package helpers

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

// Bytes generates n random bytes
func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

const (
	DIGITS_ONLY                    = "0123456789"
	LETTERS_ONLY                   = "abcdefghijklmnopqrstuvwxyz"
	LETTERS_AND_NUMBERS            = "abcdefghijklmnopqrstuvwxyz0123456789"
	MIXED_CASE_LETTERS             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	MIXED_CASE_LETTERS_AND_NUMBERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// String generates a random string using only letters provided in the letters parameter
func String(n int, letters string) string {
	var letterRunes []rune
	if letters == "" {
		letterRunes = []rune(LETTERS_ONLY)
	} else {
		letterRunes = []rune(letters)
	}

	var bb bytes.Buffer
	bb.Grow(n)
	l := uint32(len(letterRunes))
	for i := 0; i < n; i++ {
		bb.WriteRune(letterRunes[binary.BigEndian.Uint32(Bytes(4))%l])
	}
	return bb.String()
}

func Code() string {
	return String(6, DIGITS_ONLY)
}
