package main

import (
	"compress/bzip2"
	"compress/gzip"
	"context"
	"io"
	"strings"

	"cloud.google.com/go/storage"
)

type MultiObjectReader struct {
	ctx     context.Context
	bucket  *storage.BucketHandle
	objects []string
	objIdx  int
	reader  io.Reader
}

func NewMultiObjectReader(ctx context.Context, bucket *storage.BucketHandle, objects []string) *MultiObjectReader {
	return &MultiObjectReader{
		ctx,
		bucket,
		objects[:],
		0,
		nil,
	}
}

// io.Reader implementation adapted primarily from Go's io.MultiReader implementation
// (https://golang.org/src/io/multi.go)
func (m *MultiObjectReader) Read(data []byte) (n int, err error) {
	// First time only
	if m.reader == nil {
		err = m.nextReader()
		if err != nil {
			return
		}
	}
	for {
		n, err = m.reader.Read(data)
		if err == io.EOF {
			if m.objIdx < len(m.objects) {
				err = m.nextReader()
				if err != nil {
					return
				}
			}
			return
		}
		if n > 0 || err != io.EOF {
			if err == io.EOF && m.objIdx < len(m.objects) {
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}

func (m *MultiObjectReader) nextReader() error {
	nextObj := m.objects[m.objIdx]
	reader, err := m.bucket.Object(nextObj).NewReader(m.ctx)
	if err != nil {
		return err
	}
	wrapped, err := maybeDecompress(nextObj, reader)
	if err != nil {
		return err
	}
	m.objIdx += 1
	m.reader = wrapped
	return nil
}

func maybeDecompress(filename string, reader io.Reader) (io.Reader, error) {
	if *noDecompress == false {
		if strings.HasSuffix(filename, ".gz") {
			return gzip.NewReader(reader)
		}
		if strings.HasSuffix(filename, ".bz2") || strings.HasSuffix(filename, ".bzip2") {
			return bzip2.NewReader(reader), nil
		}
	}
	return reader, nil
}
