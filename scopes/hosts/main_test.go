package hosts

import "testing"

func TestCreateConnectionString(t *testing.T) {
	h := &host{}

	got := h.createConnectionString(hostHSSH{User: "deploy", Hostname: "srv.example.com"})
	if got != "deploy@srv.example.com" {
		t.Errorf("got %q", got)
	}
}

func TestCreatePortString(t *testing.T) {
	h := &host{}

	// hssh leaves the port at zero when the host config does not set one
	if got := h.createPortString(hostHSSH{}); got != "22" {
		t.Errorf("an unset port should fall back to 22, got %q", got)
	}

	if got := h.createPortString(hostHSSH{Port: 2222}); got != "2222" {
		t.Errorf("got %q", got)
	}
}
