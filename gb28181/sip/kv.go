package sip

import (
	"errors"
	"strings"
)

var (
	errEmptyKey = errors.New("empty key")
)

// KV 是一个 key-value 的结构。
// 允许单个 key 的形式
type KV struct {
	Key   string
	Value string
}

// FormatTo 格式化 k=v 到 writer 中。
func (kv *KV) FormatTo(writer Writer) error {
	return FormatKVTo2(writer, kv.Key, kv.Value)
}

// Parse 从 line 解析数据。
func (kv *KV) Parse(line string) error {
	// 找到第一个 '='
	i := strings.IndexByte(line, '=')
	if i < 0 {
		kv.Key = line
		kv.Value = ""
		return nil
	}
	// 赋值
	kv.Key = strings.TrimSpace(line[:i])
	kv.Value = strings.TrimSpace(line[i+1:])
	// key 为空
	if kv.Key == "" {
		return errEmptyKey
	}
	return nil
}

// ParseKV 从 line 解析 k=v 或者 k="v" 并返回。
func ParseKV(line string, split byte) []KV {
	var kvs []KV
	for line != "" {
		// 找到第一个 '='
		i := strings.IndexByte(line, '=')
		if i < 0 {
			kvs = append(kvs, KV{Key: line})
			break
		}
		kv := KV{Key: strings.TrimSpace(line[:i])}
		line = strings.TrimSpace(line[i+1:])
		if line == "" {
			kvs = append(kvs, kv)
			break
		}
		// 是否有双引号 '"'
		if line[0] == '"' {
			// 找到下一个 '"'
			i = strings.IndexByte(line[1:], '"')
			if i < 0 {
				kv.Value = line
				kvs = append(kvs, kv)
				break
			}
			kv.Value = line[1 : i+1]
			kvs = append(kvs, kv)
			line = line[i+2:]
			continue
		}
		// 找到分隔符
		i = strings.IndexByte(line, split)
		if i < 0 {
			kvs = append(kvs, KV{Key: line})
			break
		}
		kv.Value = line[:i]
		kvs = append(kvs, kv)
		line = line[i+1:]
	}
	return kvs
}

// FormatKVTo 格式化 key="value" 到 writer
func FormatKVTo(writer Writer, key, value string, left, right byte) error {
	_, err := writer.WriteString(key)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("=")
	if err != nil {
		return err
	}
	err = writer.WriteByte(left)
	if err != nil {
		return err
	}
	_, err = writer.WriteString(value)
	if err != nil {
		return err
	}
	err = writer.WriteByte(right)
	return err
}

// FormatKVTo2 格式化 key=value 到 writer
func FormatKVTo2(writer Writer, key, value string) error {
	_, err := writer.WriteString(key)
	if err != nil {
		return err
	}
	_, err = writer.WriteString("=")
	if err != nil {
		return err
	}
	_, err = writer.WriteString(value)
	return err
}
