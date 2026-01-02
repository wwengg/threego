//go:build !windows
// +build !windows

package utils

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
)

type BrotliCompressor struct{}

func (c BrotliCompressor) Zip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, 5)
	if _, err := w.Write(data); err != nil {
		w.Close()
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c BrotliCompressor) Unzip(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}
