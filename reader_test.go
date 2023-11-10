package rreader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestReade_Seek(t *testing.T) {
	rs := bytes.NewReader([]byte(
		"123456789qwertyuiopasdfghjklzxcvbnm",
	))
	var k2 = make([]byte, 2)
	rs.Read(k2)

	rr := NewReader(rs)

	rr2 := bytes.NewReader([]byte(
		reverseStr("3456789qwertyuiopasdfghjklzxcvbnm"),
	))

	seek, err := rr.Seek(-25, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	seek2, err := rr2.Seek(-25, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	if seek2 != seek {
		t.Fatalf("返回错误不一致 %v,%v", seek, seek2)
	}

	for {
		var k = make([]byte, 3)
		var k22 = make([]byte, 3)

		n, er := rr.Read(k)
		n2, er2 := rr2.Read(k22)
		if strErr(er) != strErr(er2) {
			t.Fatalf("返回错误不一致 %v,%v", er, er2)
		}
		if n != n2 {
			t.Fatalf("返回n不一致 %v,%v", n, n2)
		}

		ks1, ks2 := string(k[:n]), string(k22[:n2])

		//fmt.Println(ks1, ks2)
		if string(k[:n]) != string(k22[:n2]) {
			t.Fatalf("读取不一致 %v,%v", ks1, ks2)
		}
		if er == io.EOF {
			break
		}
	}

	return

}

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

//func ReverseRead(name string, lineNum uint) ([]string, error) {
//	return nil, nil
//}

func strErr(err error) string {
	if err == nil {
		return "nil"
	}
	return err.Error()
}
