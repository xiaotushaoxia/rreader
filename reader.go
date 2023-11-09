package rreader

import (
	"fmt"
	"io"
)

func NewReader(rs io.ReadSeeker) io.Reader {
	return &reader{
		rs:    rs,
		total: -1,
	}
}

type reader struct {
	rs    io.ReadSeeker
	read  int64
	total int64
}

func (r *reader) Read(bs []byte) (n int, err error) {
	if r.total < 0 {
		r.total, err = r.rs.Seek(0, io.SeekEnd)
		if err != nil {
			return
		}
	}
	if r.read >= r.total {
		return 0, io.EOF
	}
	if len(bs) == 0 && r.read == 0 {
		// 如果 bs是空的 且 没都读过字节 就会调用Seek(0, io.SeekEnd) 会导致Read返回EOF
		// 如果 bs是空的 但 有读过字节 那么Read自然会返回正确的东西
		return 0, nil
	}

	p := -r.read - int64(len(bs))
	if p < -r.total {
		bs = bs[:r.total-r.read]
		p = -r.total
	}

	_, err = r.rs.Seek(p, io.SeekEnd)
	if err != nil {
		return 0, failedTo("seek", err)
	}

	n, err = r.rs.Read(bs)
	if n > 0 {
		r.read += int64(n)
		reverse(bs[:n])
	}
	return
}

func failedTo(opt string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to "+opt+" %w", err)
}

func reverse[T any](is []T) []T {
	i, j := 0, len(is)-1
	for i < j {
		is[i], is[j] = is[j], is[i]
		i++
		j--
	}
	return is
}
