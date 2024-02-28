package cachestorage

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
)

const cacheIdentifier = ":cache@identifier:"

type diskCacheStorage struct {
	dest string
}

func NewDiskCacheStorage(destination string) CacheStorage {
	return &diskCacheStorage{
		dest: destination,
	}
}

func (d *diskCacheStorage) Push(name string, fsys CacheableFS) error {
	f, err := os.OpenFile(path.Join(d.dest, name+".zip"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open zip file: %w", err)
	}

	w := zip.NewWriter(f)
	defer func(w *zip.Writer) {
		err := w.Close()
		if err != nil {
			slog.Error("failed to close zip writer", slog.String("error", err.Error()))
		}
	}(w)

	// FIXME: Go 1.22 - archive/zip.AddFS
	err = fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("info: %w", err)
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("file info header: %w", err)
		}
		h.Name = name
		h.Method = zip.Deflate
		fw, err := w.CreateHeader(h)
		if err != nil {
			return fmt.Errorf("create header: %w", err)
		}

		f, err := fsys.Open(name)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer func(f fs.File) {
			err := f.Close()
			if err != nil {
				slog.Error("failed to close file", slog.String("error", err.Error()))
			}
		}(f)

		_, err = io.Copy(fw, f)
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err := w.SetComment(cacheIdentifier); err != nil {
		return err
	}

	return nil
}

type cachedDiskFile struct {
	z *zip.ReadCloser
}

func (c *cachedDiskFile) Open(name string) (fs.File, error) {
	return c.z.Open(name)
}

func (c *cachedDiskFile) Extract(dst string) error {
	for _, f := range c.z.File {
		dstName := path.Join(dst, f.Name)
		slog.Debug("extracting file", slog.String("file", dstName))

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(dstName, 0755)
			if err != nil {
				return fmt.Errorf("mkdir (for dir): %w", err)
			}
			continue
		}

		if err := os.MkdirAll(path.Dir(dstName), 0755); err != nil {
			return fmt.Errorf("mkdir (for file): %w", err)
		}

		r, err := f.Open()
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}

		w, err := os.Create(dstName)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		err = os.Chmod(dstName, f.Mode())
		if err != nil {
			return fmt.Errorf("chmod: %w", err)
		}

		_, err = io.Copy(w, r)
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}

		_ = w.Close()
	}

	return nil
}

func (c *cachedDiskFile) Close() error {
	return c.z.Close()
}

// Underlying returns the underlying zip reader.
//
// It is for test only and should not be used in production code.
func (c *cachedDiskFile) Underlying() *zip.ReadCloser {
	return c.z
}

func (d *diskCacheStorage) Pull(name string) (CacheImage, error) {
	r, err := zip.OpenReader(path.Join(d.dest, name+".zip"))
	if err != nil || r.Comment != cacheIdentifier {
		// invalidate it if it's not a valid cache
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error("failed to open zip reader", slog.String("error", err.Error()))
			_ = d.Invalidate(name)
		}

		return nil, fmt.Errorf("open zip reader: %w", err)
	}

	return &cachedDiskFile{
		z: r,
	}, nil
}

func (d *diskCacheStorage) Invalidate(name string) error {
	return os.Remove(path.Join(d.dest, name+".zip"))
}

var _ CacheStorage = (*diskCacheStorage)(nil)
