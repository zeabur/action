package cachestorage_test

import (
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"testing"

	"github.com/psanford/memfs"
	"github.com/zeabur/action/cachestorage"
	"gotest.tools/v3/assert"
)

func TestMain(m *testing.M) {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(h))

	os.Exit(m.Run())
}

func TestDiskStorage_FullFlow(t *testing.T) {
	dir, err := os.MkdirTemp("", "zbaction-testcase-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dir)

	tmpfs := memfs.New()
	if err := tmpfs.WriteFile("test-file", []byte("test-value"), 0755); err != nil {
		t.Error(err)
	}
	if err := tmpfs.MkdirAll("dir-1", 0755); err != nil {
		t.Error(err)
	}
	if err := tmpfs.WriteFile("dir-1/test-file2", []byte("test-value2"), 0755); err != nil {
		t.Error(err)
	}

	storage := cachestorage.NewDiskCacheStorage(dir)
	err = storage.Push("test-disk-storage", tmpfs)
	if err != nil {
		t.Fatal(err)
	}

	ci, err := storage.Pull("test-disk-storage")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("pull.read-file", func(t *testing.T) {
		f, err := ci.Open("test-file")
		if err != nil {
			t.Fatal(err)
		}
		defer func(f fs.File) {
			_ = f.Close()
		}(f)

		b, err := io.ReadAll(f)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "test-value", string(b))
	})

	t.Run("pull.read-file-nested", func(t *testing.T) {
		f, err := ci.Open("dir-1/test-file2")
		if err != nil {
			t.Fatal(err)
			return
		}
		defer func(f fs.File) {
			_ = f.Close()
		}(f)

		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatal(err)
			return
		}

		assert.Equal(t, "test-value2", string(b))
	})

	t.Run("pull.extract", func(t *testing.T) {
		testtreedir, err := os.MkdirTemp("", "zbaction-test-tree-*")
		if err != nil {
			t.Fatal(err)
		}
		defer func(path string) {
			_ = os.RemoveAll(path)
		}(testtreedir)

		err = ci.Extract(dir)
		if err != nil {
			t.Fatal(err)
		}

		c, err := os.ReadFile(path.Join(dir, "test-file"))
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, "test-value", string(c))
		}

		c, err = os.ReadFile(path.Join(dir, "dir-1/test-file2"))
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, "test-value2", string(c))
		}
	})
}
