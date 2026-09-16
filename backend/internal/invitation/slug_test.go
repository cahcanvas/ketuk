package invitation

import "testing"

func TestSlugify(t *testing.T) {
	if got := slugify("Akad & Resepsi Budi"); got != "akad-resepsi-budi" {
		t.Fatalf("got %q", got)
	}
	if got := slugify("   "); got != "undangan" {
		t.Fatalf("got %q", got)
	}
}
