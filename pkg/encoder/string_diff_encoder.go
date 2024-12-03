package encoder

import (
	"errors"
	"fmt"
	"github.com/compression-algorithm-research-lab/go-string-delta/pkg/diff"
)

// TODO 2024-12-04 01:27:11 待Review、Test

// TODO 2024-12-04 00:45:54 是否要加一个可选的类型头？这样子的话拿到一个字节数组的时候能够区分到底是什么类型的编码

type DefaultStringDiffEncoder struct {
}

var _ StringDiffEncoder = &DefaultStringDiffEncoder{}

func (x *DefaultStringDiffEncoder) Encode(d *diff.StringDiff) ([]byte, error) {
	bytes := make([]byte, 0)
	switch d.ChangeType {
	case diff.ChangeTypeNone:
		bytes = append(bytes, byte(diff.ChangeTypeNone))
		break
	case diff.ChangeTypeInsert:
		bytes = append(bytes, byte(diff.ChangeTypeInsert))
		bytes = x.appendUint64(bytes, d.BeginOffset)
		bytes = x.appendString(bytes, d.Content)
		break
	case diff.ChangeTypeDelete:
		bytes = append(bytes, byte(diff.ChangeTypeDelete))
		bytes = x.appendUint64(bytes, d.BeginOffset)
		bytes = x.appendUint64(bytes, d.EndOffset)
		break
	case diff.ChangeTypeReplace:
		bytes = append(bytes, byte(diff.ChangeTypeReplace))
		bytes = x.appendUint64(bytes, d.BeginOffset)
		bytes = x.appendUint64(bytes, d.EndOffset)
		bytes = x.appendString(bytes, d.Content)
	}
	return bytes, nil
}

func (x *DefaultStringDiffEncoder) EncodeSlice(diffSlice []*diff.StringDiff) ([]byte, error) {

	// TODO 2024-12-04 00:34:42 是否需要先排序？然后再压缩，这样相同类型的可以节省一个type

	bytes := make([]byte, 0)
	for _, d := range diffSlice {
		encodeBytes, err := x.Encode(d)
		if err != nil {
			return nil, err
		}
		bytes = x.appendUint64(bytes, uint64(len(encodeBytes)))
		bytes = append(bytes, encodeBytes...)
	}
	return bytes, nil
}

func (x *DefaultStringDiffEncoder) EncodeSequence(first string, diffSequence [][]*diff.StringDiff) ([]byte, error) {
	bytes := make([]byte, 0)

	// 第一个字节放进去
	bytes = x.appendString(bytes, first)

	// 变化部分
	for _, diffSlice := range diffSequence {
		encodeBytes, err := x.EncodeSlice(diffSlice)
		if err != nil {
			return nil, err
		}
		bytes = x.appendBytes(bytes, encodeBytes)
	}
	return bytes, nil
}

func (x *DefaultStringDiffEncoder) Decode(bytes []byte) (*diff.StringDiff, error) {

	// 第一个字节是变化的类型
	if len(bytes) < 1 {
		return nil, fmt.Errorf("bytes length error")
	}

	d := &diff.StringDiff{}
	d.ChangeType = diff.ChangeType(bytes[0])
	index := 1
	var err error
	var n uint64
	var s string
	switch d.ChangeType {
	case diff.ChangeTypeNone:
		break
	case diff.ChangeTypeInsert:

		// BeginOffset
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return nil, err
		}
		d.BeginOffset = n

		// Content
		index, s, err = x.readString(bytes, index)
		if err != nil {
			return nil, err
		}
		d.Content = s

		break
	case diff.ChangeTypeDelete:

		// BeginOffset
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return nil, err
		}
		d.BeginOffset = n

		// EndOffset
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return nil, err
		}
		d.EndOffset = n

		break
	case diff.ChangeTypeReplace:

		// BeginOffset
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return nil, err
		}
		d.BeginOffset = n

		// EndOffset
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return nil, err
		}
		d.EndOffset = n

		// Content
		index, s, err = x.readString(bytes, index)
		if err != nil {
			return nil, err
		}
		d.Content = s

		break
	}

	return d, nil
}

func (x *DefaultStringDiffEncoder) DecodeSlice(bytes []byte) ([]*diff.StringDiff, error) {
	var diffSlice []*diff.StringDiff
	index := 0
	var err error
	var n uint64
	for index < len(bytes) {
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return nil, err
		}
		diffBytes := bytes[index : index+int(n)]
		d, err := x.Decode(diffBytes)
		if err != nil {
			return nil, err
		}
		diffSlice = append(diffSlice, d)
		index += int(n)
	}
	return diffSlice, nil
}

func (x *DefaultStringDiffEncoder) DecodeSequence(bytes []byte) (first string, diffSequence [][]*diff.StringDiff, err error) {
	index := 0
	var n uint64
	index, first, err = x.readString(bytes, index)
	if err != nil {
		return "", nil, err
	}
	for index < len(bytes) {
		index, n, err = x.readUint64(bytes, index)
		if err != nil {
			return "", nil, err
		}
		diffSlice, err := x.DecodeSlice(bytes[index : index+int(n)])
		if err != nil {
			return "", nil, err
		}
		diffSequence = append(diffSequence, diffSlice)
		index += int(n)
	}
	return first, diffSequence, nil
}

// ------------------------------------------------ ---------------------------------------------------------------------

// 往字节数组中追加一个字符串
func (x *DefaultStringDiffEncoder) appendString(bytes []byte, s string) []byte {
	sBytes := []byte(s)
	// 写入字符串字节长度
	bytes = x.appendUint64(bytes, uint64(len(sBytes)))
	// 写入字符串字节数组
	bytes = append(bytes, sBytes...)
	return bytes
}

// 从字节数组的给定位置读取一个字符串
func (x *DefaultStringDiffEncoder) readString(bytes []byte, index int) (int, string, error) {
	index, n, err := x.readUint64(bytes, index)
	if err != nil {
		return 0, "", err
	}
	if len(bytes) < index+int(n) {
		return 0, "", fmt.Errorf("string length error")
	}
	sBytes := bytes[index : index+int(n)]
	return index + int(n), string(sBytes), nil
}

// ------------------------------------------------ ---------------------------------------------------------------------

// 往字节数组中写入一个uint64
func (x *DefaultStringDiffEncoder) appendUint64(bytes []byte, n uint64) []byte {
	bytes = append(bytes, x.uint64ToBytes(n)...)
	return bytes
}

// 从字节数组中读取一个uint64值
func (x *DefaultStringDiffEncoder) readUint64(bytes []byte, index int) (int, uint64, error) {
	if len(bytes) < index+8 {
		return 0, 0, errors.New("string length error")
	}
	toUint64 := x.bytesToUint64(bytes[index : index+8])
	return index + 8, toUint64, nil
}

// ------------------------------------------------ ---------------------------------------------------------------------

// 往字节数组中写入一个uint64
func (x *DefaultStringDiffEncoder) appendBytes(bytes []byte, bytesForAppend []byte) []byte {
	bytes = x.appendUint64(bytes, uint64(len(bytesForAppend)))
	bytes = append(bytes, bytesForAppend...)
	return bytes
}

// 从字节数组中读取一个uint64值
func (x *DefaultStringDiffEncoder) readBytes(bytes []byte, index int) (int, []byte, error) {
	index, n, err := x.readUint64(bytes, index)
	if err != nil {
		return 0, nil, err
	}
	if len(bytes) < index+int(n) {
		return 0, nil, errors.New("string length error")
	}

	return index + int(n), bytes[index : index+int(n)], nil
}

// ------------------------------------------------ ---------------------------------------------------------------------

// TODO 2024-12-03 23:48:51 长度编码自适应
// UintToBytes 将无符号整数转换为字节切片
func (x *DefaultStringDiffEncoder) uint64ToBytes(value uint64) []byte {
	var byteValue [8]byte
	byteValue[7] = byte(value)
	value >>= 8
	byteValue[6] = byte(value)
	value >>= 8
	byteValue[5] = byte(value)
	value >>= 8
	byteValue[4] = byte(value)
	value >>= 8
	byteValue[3] = byte(value)
	value >>= 8
	byteValue[2] = byte(value)
	value >>= 8
	byteValue[1] = byte(value)
	value >>= 8
	byteValue[0] = byte(value)

	return byteValue[:]
}

func (x *DefaultStringDiffEncoder) bytesToUint64(bytes []byte) uint64 {
	// TODO 2024-12-04 01:03:44
	return 0
}

// ------------------------------------------------ ---------------------------------------------------------------------
