package rreader

import (
	"fmt"
	"io"
)

var (
	ErrNegativePosition = fmt.Errorf("negative position")
	ErrInvalidWhence    = fmt.Errorf("invalid whence")
)

func NewReader(rs io.ReadSeeker) io.ReadSeeker {
	return &reader{
		rs: rs,
	}
}

type reader struct {
	rs    io.ReadSeeker
	read  int64
	total int64

	// 保证getTotal只调用一次 不考虑并发 不用sync.Once
	getTotalCalled bool
	getTotalErr    error
}

func (r *reader) Read(bs []byte) (n int, err error) {
	err = r.getTotal()
	if err != nil {
		return
	}
	if r.read >= r.total {
		return 0, io.EOF
	}
	if len(bs) == 0 && r.read == 0 {
		// 如果 bs是空的 且 没都读过字节 就会调用Seek(0, io.SeekEnd) 会导致Read返回EOF
		// 如果 bs是空的 但 有读过字节 那么Read自然会返回正确的东西
		return 0, nil
	}

	//p := max1(-r.read-int64(len(bs)), -r.total)
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

func (r *reader) getTotal() error {
	if r.getTotalCalled {
		return r.getTotalErr
	}
	r.getTotalCalled = true

	//原本的实现 如果r.rs读到过就会重复读那部分
	// 而且之前用的 r.total, err := r.rs.Seek(0, io.SeekEnd) 会有问题
	//temp, err := r.rs.Seek(0, io.SeekEnd)
	//if err != nil {
	//	return err
	//}
	//r.total = temp
	//return nil

	// 修改后：r.rs已经读过的那部分不再读
	cur, err := r.rs.Seek(0, io.SeekCurrent)
	if err != nil {
		r.getTotalErr = err
		return err
	}
	sum, err := r.rs.Seek(0, io.SeekEnd)
	if err != nil {
		r.getTotalErr = err
		return err
	}
	r.total = sum - cur
	return nil
}

func (r *reader) Seek(offset int64, whence int) (int64, error) {
	err := r.getTotal()
	if err != nil {
		return 0, err
	}
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.read + offset
	case io.SeekEnd:
		abs = r.total + offset
	default:
		return 0, ErrInvalidWhence
	}
	if abs < 0 {
		return 0, ErrNegativePosition
	}
	r.read = abs
	return abs, nil
}

func failedTo(opt string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to "+opt+" %w", err)
}

func reverse[T any](is []T) []T {
	//i, j := 0, len(is)-1
	//for i < j {
	//	is[i], is[j] = is[j], is[i]
	//	i++
	//	j--
	//}
	for i, j := 0, len(is)-1; i < j; i, j = i+1, j-1 {
		is[i], is[j] = is[j], is[i]
	}
	return is
}
