package every

import (
	"errors"
	"fmt"
)

type Buffer struct {
	data string
}

func (b *Buffer) Read() string {
	fmt.Println("Read: ", b.data)
	return b.data
}

func (b *Buffer) Write(s string) {
	b.data = s
	fmt.Println("Write: ", s)
}

type Reader interface {
	Read() string
}

type Writer interface {
	Write(s string)
}

type ReadWriter interface {
	Reader
	Writer
}

func CopyData(r Reader, w Writer) {
	data := r.Read()
	w.Write(data)
}

type Closer interface {
	Close() error
}

func (b *Buffer) Close() error {
	if b.data == "" {
		return errors.New("пустая строка")
	}
	return nil
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}

func embedding() {
	var r ReadWriteCloser = &Buffer{data: "data"}
	r.Write("hello")
	data := r.Read()
	fmt.Println(data)
	err := r.Close()
	if err != nil {
		fmt.Println("ошибка: ", err)
	}
}
