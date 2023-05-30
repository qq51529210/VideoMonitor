package util

import (
	"bufio"
	"io"
	"strings"
)

// LineReader 用于读取 lf 的每一行数据
type LineReader struct {
	// 数据源
	scaner *bufio.Scanner
}

// ReadLine读取一行数据
func (r *LineReader) ReadLine() (string, error) {
	var line strings.Builder
	for r.scaner.Scan() {
		b := r.scaner.Bytes()
		if len(b) < 1 {
			continue
		}
		line.Write(b)
		break
	}
	if err := r.scaner.Err(); err != nil && err != io.EOF {
		return "", err
	}
	return line.String(), nil
}

// NewLineReader 返回新的 LineReader
func NewLineReader(reader io.Reader) *LineReader {
	r := new(LineReader)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	r.scaner = scanner
	return r
}
