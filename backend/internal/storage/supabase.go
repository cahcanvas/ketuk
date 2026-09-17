package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Supabase talks to a Supabase Storage bucket over its REST API — no SDK,
// same convention as internal/pay's Duitku adapter. The bucket must be
// public: invitation media (gallery/music) has no access control of its own,
// it is meant to be viewable by anyone with the RSVP link.
type Supabase struct {
	BaseURL    string
	Bucket     string
	ServiceKey string
	HTTP       *http.Client
}

func NewSupabase(baseURL, bucket, serviceKey string) *Supabase {
	return &Supabase{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Bucket:     bucket,
		ServiceKey: serviceKey,
		HTTP:       &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *Supabase) objectURL(path string) string {
	return fmt.Sprintf("%s/storage/v1/object/%s/%s", s.BaseURL, s.Bucket, path)
}

func (s *Supabase) Put(ctx context.Context, path, contentType string, src io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.objectURL(path), src)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.ServiceKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")
	resp, err := s.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase storage upload failed: %s: %s", resp.Status, body)
	}
	return nil
}

func (s *Supabase) Delete(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, s.objectURL(path), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.ServiceKey)
	resp, err := s.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase storage delete failed: %s: %s", resp.Status, body)
	}
	return nil
}

// URL points at the bucket's public object endpoint, so callers redirect
// instead of proxying bytes through our own function.
func (s *Supabase) URL(path string) string {
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.BaseURL, s.Bucket, path)
}

// Open is never called: URL is always non-empty for this store.
func (s *Supabase) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("supabase store: Open is not supported, use URL")
}
