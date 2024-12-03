package encoder

import "github.com/compression-algorithm-research-lab/go-string-delta/pkg/diff"

// StringDiffEncoder 编码器，将其编码为字节数组或者从字节数组解码
type StringDiffEncoder interface {

	// Encode 将单个的diff信息编码为字节数组
	Encode(diff *diff.StringDiff) ([]byte, error)

	// EncodeSlice 将多个diff信息编码为字节数组
	EncodeSlice(diffSlice []*diff.StringDiff) ([]byte, error)

	// EncodeSequence 一整个序列
	EncodeSequence(first string, diffSequence [][]*diff.StringDiff) ([]byte, error)

	// Decode 用于从编码压缩的字节数组还原出原始的变化情况
	Decode(data []byte) (*diff.StringDiff, error)

	// DecodeSlice 用于从编码压缩的字节数组还原出原始的变化数组
	DecodeSlice(data []byte) ([]*diff.StringDiff, error)

	// DecodeSequence 一整个序列
	DecodeSequence(data []byte) (string, [][]*diff.StringDiff, error)
}
