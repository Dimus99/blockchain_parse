package main

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

func Reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
func Reverse_binary(s []byte) (result [32]byte) {
	result = [32]byte{}
	for i, v := range s {
		result[32-i-1] = v
	}
	return
}

//func B2S(bs []uint8) string {
//	b := make([]byte, len(bs))
//	for i, v := range bs {
//		b[i] = v
//	}
//	return string(b)
//}

func BE_TO_LE(num string) string {
	num = Reverse(num)
	res := []uint8{}
	if len(num)%2 != 0 {
		panic(num)
	}
	for i := 0; i < utf8.RuneCountInString(num); i += 2 {
		res = append(res, num[i+1])
		res = append(res, num[i])
	}

	return string(res)
}

func zfill(s string) string {
	for utf8.RuneCountInString(s) < 8 {
		s = "0" + s
	}
	return s
}

func hexlify(s string) string {
	ui, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		panic(err)
	}
	panic(fmt.Sprintf("%x", ui))
	return fmt.Sprintf("%x", ui)
}
