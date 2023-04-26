package tools

import (
	"fmt"
	"strconv"
	"testing"
)

func TestFillRandKey(t *testing.T) {
	s := FillRandKey(10)
	fmt.Println(s)
}
func TestStringTransBinary(t *testing.T) {
	s := fmt.Sprintf("%b", 'A')
	//1000001
	fmt.Println(s)
	k, _ := strconv.ParseInt("1000001", 2, 64)

	fmt.Println(byte(k))
}

func TestCalRandKey(t *testing.T) {
	oldRandKey := FillRandKey(36)
	wordList := GetWordList(oldRandKey, 9)
	newRandKey := calRandKey(wordList, 9)

	fmt.Println(oldRandKey == newRandKey)
}
