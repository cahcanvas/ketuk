package invitation

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"ketuk.id/api/internal/apierr"
)

type Media struct {
	ID           uuid.UUID `json:"id"`
	Kind         string    `json:"kind"`
	OriginalName string    `json:"original_name"`
	ContentType  string    `json:"content_type"`
	URL          string    `json:"url"`
	Caption      string    `json:"caption"`
	SortOrder    int       `json:"sort_order"`
	relPath      string
}

func (s *Service) AddGallery(ctx context.Context, owner, invitationID uuid.UUID, original, contentType, caption string, src io.Reader) (Media, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Media{}, err
	}
	if !isImage(contentType) {
		return Media{}, apierr.BadRequest("invalid_input", "gallery must be an image")
	}
	return s.saveMedia(ctx, invitationID, "gallery", original, contentType, caption, src)
}

func (s *Service) SetMusic(ctx context.Context, owner, invitationID uuid.UUID, original, contentType string, src io.Reader) (Media, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Media{}, err
	}
	if !isAudio(contentType) {
		return Media{}, apierr.BadRequest("invalid_input", "music must be an audio file")
	}
	_, err := s.pool.Exec(ctx, `DELETE FROM invitation_media WHERE invitation_id=$1 AND kind='music'`, invitationID)
	if err != nil {
		return Media{}, err
	}
	return s.saveMedia(ctx, invitationID, "music", original, contentType, "", src)
}

func (s *Service) ListGallery(ctx context.Context, owner, invitationID uuid.UUID) ([]Media, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	gallery, _, err := s.listMedia(ctx, invitationID)
	return gallery, err
}

func (s *Service) GetMusic(ctx context.Context, owner, invitationID uuid.UUID) (*Media, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return nil, err
	}
	_, music, err := s.listMedia(ctx, invitationID)
	return music, err
}

// GalleryPatch only carries the caption; the stored file, kind, and sort order stay put.
type GalleryPatch struct {
	Caption *string `json:"caption"`
}

func (s *Service) PatchGallery(ctx context.Context, owner, invitationID, mediaID uuid.UUID, in GalleryPatch) (Media, error) {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return Media{}, err
	}
	var m Media
	err := s.pool.QueryRow(ctx, `
		SELECT id, kind, original_name, content_type, rel_path, caption, sort_order
		FROM invitation_media WHERE id=$1 AND invitation_id=$2 AND kind='gallery'`, mediaID, invitationID).
		Scan(&m.ID, &m.Kind, &m.OriginalName, &m.ContentType, &m.relPath, &m.Caption, &m.SortOrder)
	if errors.Is(err, pgx.ErrNoRows) {
		return Media{}, apierr.NotFound("media")
	}
	if err != nil {
		return Media{}, err
	}
	m.URL = "/v1/public/media/" + m.ID.String()
	m = applyGalleryCaption(m, in)
	_, err = s.pool.Exec(ctx, `
		UPDATE invitation_media SET caption=$3
		WHERE id=$1 AND invitation_id=$2 AND kind='gallery'`,
		m.ID, invitationID, m.Caption)
	return m, err
}

// applyGalleryCaption trims a sent caption; an empty string clears it, nil leaves it alone.
func applyGalleryCaption(m Media, patch GalleryPatch) Media {
	if patch.Caption != nil {
		m.Caption = strings.TrimSpace(*patch.Caption)
	}
	return m
}

func (s *Service) DeleteMedia(ctx context.Context, owner, invitationID, mediaID uuid.UUID) error {
	if _, err := s.Get(ctx, owner, invitationID); err != nil {
		return err
	}
	var rel string
	err := s.pool.QueryRow(ctx, `SELECT rel_path FROM invitation_media WHERE id=$1 AND invitation_id=$2`, mediaID, invitationID).Scan(&rel)
	if err != nil {
		return apierr.NotFound("media")
	}
	_, err = s.pool.Exec(ctx, `DELETE FROM invitation_media WHERE id=$1`, mediaID)
	if err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(s.storageDir, rel))
	return nil
}

func (s *Service) OpenMedia(ctx context.Context, mediaID uuid.UUID) (absPath, contentType, original string, err error) {
	var rel string
	err = s.pool.QueryRow(ctx, `
		SELECT m.rel_path, m.content_type, m.original_name
		FROM invitation_media m WHERE m.id=$1`, mediaID).
		Scan(&rel, &contentType, &original)
	if err != nil {
		return "", "", "", apierr.NotFound("media")
	}
	return filepath.Join(s.storageDir, rel), contentType, original, nil
}

func (s *Service) saveMedia(ctx context.Context, invitationID uuid.UUID, kind, original, contentType, caption string, src io.Reader) (Media, error) {
	id := uuid.New()
	ext := extFor(original, contentType)
	rel := filepath.ToSlash(filepath.Join("invitations", invitationID.String(), kind, id.String()+ext))
	abs := filepath.Join(s.storageDir, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return Media{}, err
	}
	f, err := os.Create(abs)
	if err != nil {
		return Media{}, err
	}
	defer f.Close()
	if _, err := io.Copy(f, src); err != nil {
		_ = os.Remove(abs)
		return Media{}, err
	}
	m := Media{
		ID: id, Kind: kind, OriginalName: original, ContentType: contentType,
		Caption: strings.TrimSpace(caption), URL: "/v1/public/media/" + id.String(),
	}
	_, err = s.pool.Exec(ctx, `
		INSERT INTO invitation_media (id, invitation_id, kind, original_name, content_type, rel_path, caption)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		m.ID, invitationID, kind, original, contentType, rel, m.Caption)
	return m, err
}

func (s *Service) listMedia(ctx context.Context, invitationID uuid.UUID) ([]Media, *Media, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, kind, original_name, content_type, rel_path, caption, sort_order
		FROM invitation_media WHERE invitation_id=$1 ORDER BY sort_order, created_at`, invitationID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	gallery := []Media{}
	var music *Media
	for rows.Next() {
		var m Media
		if err := rows.Scan(&m.ID, &m.Kind, &m.OriginalName, &m.ContentType, &m.relPath, &m.Caption, &m.SortOrder); err != nil {
			return nil, nil, err
		}
		m.URL = "/v1/public/media/" + m.ID.String()
		if m.Kind == "music" {
			cp := m
			music = &cp
			continue
		}
		gallery = append(gallery, m)
	}
	return gallery, music, rows.Err()
}

func isImage(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "image/")
}

func isAudio(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "audio/")
}

func extFor(name, contentType string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".jpeg" {
		return ".jpg"
	}
	if ext != "" && len(ext) <= 8 {
		return ext
	}
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "audio/mpeg":
		return ".mp3"
	case "audio/wav":
		return ".wav"
	default:
		return ".bin"
	}
}
