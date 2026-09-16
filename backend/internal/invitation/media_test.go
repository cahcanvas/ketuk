package invitation

import (
	"testing"

	"github.com/google/uuid"
)

func TestApplyGalleryCaption(t *testing.T) {
	id := uuid.New()
	current := Media{
		ID: id, Kind: "gallery", OriginalName: "akad.jpg", ContentType: "image/jpeg",
		URL: "/v1/public/media/" + id.String(), Caption: "Akad nikah", SortOrder: 3,
	}

	if got := applyGalleryCaption(current, GalleryPatch{}); got != current {
		t.Fatalf("nil caption changed media: %+v", got)
	}

	caption := "  Sesi akad  "
	got := applyGalleryCaption(current, GalleryPatch{Caption: &caption})
	if got.Caption != "Sesi akad" {
		t.Fatalf("got caption %q", got.Caption)
	}
	if got.URL != current.URL || got.Kind != current.Kind || got.SortOrder != current.SortOrder || got.ContentType != current.ContentType {
		t.Fatalf("patch touched more than caption: %+v", got)
	}

	blank := "   "
	if cleared := applyGalleryCaption(current, GalleryPatch{Caption: &blank}); cleared.Caption != "" {
		t.Fatalf("expected caption cleared, got %q", cleared.Caption)
	}
}
