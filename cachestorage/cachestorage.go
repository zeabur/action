// Package cachestorage provides a simple storage interface based on fs.FS.
package cachestorage

import (
	"io"
	"io/fs"
)

type CacheableFS interface {
	fs.FS
}

type CacheImage interface {
	fs.FS
	io.Closer

	Extract(dst string) error
}

type CacheStorage interface {
	Push(name string, fs CacheableFS) error
	Pull(name string) (CacheImage, error)
	Invalidate(name string) error
}
