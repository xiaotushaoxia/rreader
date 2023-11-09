package rreader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestReader_Read(t *testing.T) {
	rs := bytes.NewReader([]byte(
		"123454-3-2-1-",
	))

	nn := bytes.NewBuffer(nil)

	newReader := NewReader(rs)
	n, err := nn.ReadFrom(newReader)
	fmt.Println(n, err)
	fmt.Println(nn.String()) // 输出-1-2-3-454321
}

func TestReader_ReadWithR(t *testing.T) {
	rs := bytes.NewReader([]byte(
		"111\r\n2222\r\n33333\r\n444444\r\n3333333\r\n22222222\r\n111111111",
	))
	sc := bufio.NewScanner(NewReader(rs))
	for sc.Scan() {
		aa := sc.Text()
		fmt.Println(aa)
	}
}

func TestReader_Found(t *testing.T) {
	rs := bytes.NewReader([]byte(
		"111\n2222\n33333\n444444\n3333333\n22222222\n111111111",
	))

	fmt.Println(firstLineIWant(rs))
	rs = bytes.NewReader([]byte(
		"111\r\n2222\r\n33333\n444444\n3333333\n22222222\n111111111",
	))
	fmt.Println(lastLineIWant(rs))
}

func firstLineIWant(file io.ReadSeeker) (target string, miss bool) {
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		target = sc.Text()
		if isThisLineIWant(target) {
			return
		}
	}
	miss = true
	return
}

func lastLineIWant(file io.ReadSeeker) (target string, miss bool) {
	sc := bufio.NewScanner(NewReader(file))
	for sc.Scan() {
		target = reverseStr(sc.Text())
		if isThisLineIWant(target) {
			return
		}
	}
	miss = true
	return
}

func isThisLineIWant(l string) bool {
	if len(l) == 0 {
		return false
	}
	return l[0] == '2'
}

func reverseStr(s string) string {
	a := []rune(s)
	reverse(a)
	return string(a)
}

func ReverseRead(name string, lineNum uint) ([]string, error) {
	return nil, nil
}
