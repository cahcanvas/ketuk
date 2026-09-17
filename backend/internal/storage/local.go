package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// Local writes media straight to disk under Dir. Used for local dev and for
// any host with a durable filesystem (Vercel is not one of them).
type Local struct {
	Dir string
}

func NewLocal(dir string) *Local {
	if dir == "" {
		dir = "./storage"
	}
	return &Local{Dir: dir}
}

func (l *Local) Put(ctx context.Context, path, contentType string, src io.Reader) error {
	abs := filepath.Join(l.Dir, path)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return err
	}
	f, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, src); err != nil {
		_ = os.Remove(abs)
		return err
	}
	return nil
}

func (l *Local) Delete(ctx context.Context, path string) error {
	err := os.Remove(filepath.Join(l.Dir, path))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// URL is always empty: local files have no public URL scheme of their own,
// so callers must stream bytes through Open instead.
func (l *Local) URL(path string) string { return "" }

func (l *Local) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(l.Dir, path))
}
