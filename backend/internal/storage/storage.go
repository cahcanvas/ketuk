// Package storage is the port media uploads (invitation gallery/music) go
// through. Local persists to disk — fine on a machine with a durable
// filesystem (dev, a VPS). Supabase persists to a Supabase Storage bucket —
// needed on serverless hosts (Vercel) where the filesystem is not durable
// across invocations.
package storage

import (
	"context"
	"io"
)

type Store interface {
	Put(ctx context.Context, path, contentType string, src io.Reader) error
	Delete(ctx context.Context, path string) error
	// URL returns an absolute, publicly fetchable URL for path, or "" if the
	// store has no such URL and the caller must stream bytes via Open.
	URL(path string) string
	Open(ctx context.Context, path string) (io.ReadCloser, error)
}
