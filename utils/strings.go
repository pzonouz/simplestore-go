package utils

import (
	"math/rand/v2"
	"strings"
)

func GetRandomNumber(min int, max int) int {
	return min + int(rand.IntN(max-min+1))
}

func CreateRandomString(l int) string {
	const chars = "abcdefghijklmnpqrstuvwxyz"
	k := len(chars)
	var sb strings.Builder
	for i := 0; i < l; i++ {
		sb.WriteByte(chars[GetRandomNumber(0, k-1)])
	}
	return sb.String()
}
func CreateRandomNumberAsString(l int) string {
	const chars = "0123456789"
	k := len(chars)
	var sb strings.Builder
	for i := 0; i < l; i++ {
		sb.WriteByte(chars[GetRandomNumber(0, k-1)])
	}
	return sb.String()
}
func CreateRandomPhoneNumberAsString() string {
	const chars = "0123456789"
	k := len(chars)
	var sb strings.Builder
	sb.Write([]byte("09"))
	for i := 0; i < 9; i++ {
		sb.WriteByte(chars[GetRandomNumber(0, k-1)])
	}
	return sb.String()
}
