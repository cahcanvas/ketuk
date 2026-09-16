package invitation

import "testing"

func TestValidateHostGuestStatus(t *testing.T) {
	for _, ok := range []string{"draft", " invited "} {
		got, err := validateHostGuestStatus(ok)
		if err != nil {
			t.Fatalf("%q: %v", ok, err)
		}
		if got != "draft" && got != "invited" {
			t.Fatalf("%q: got %q", ok, got)
		}
	}
	// RSVP answers belong to the guest, not the host.
	for _, bad := range []string{"", "opened", "accepted", "declined", "Invited"} {
		if _, err := validateHostGuestStatus(bad); err == nil {
			t.Fatalf("expected %q rejected", bad)
		}
	}
}
