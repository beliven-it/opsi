package onepassword

import (
	"errors"
	"strings"
	"testing"
)

const usersJSON = `[
  {"id":"1","name":"Alice","email":"alice@example.com","state":"SUSPENDED"},
  {"id":"2","name":"Bob","email":"bob@example.com","state":"ACTIVE"},
  {"id":"3","name":"Carol","email":"carol@example.com","state":"SUSPENDED"}
]`

// fakeOp replaces the op binary and records the deleted user IDs.
func fakeOp(t *testing.T, failOn string) *[]string {
	t.Helper()
	deleted := &[]string{}
	original := execOp
	t.Cleanup(func() { execOp = original })

	execOp = func(args ...string) ([]byte, error) {
		switch {
		case len(args) >= 2 && args[0] == "user" && args[1] == "list":
			return []byte(usersJSON), nil
		case len(args) == 3 && args[0] == "user" && args[1] == "delete":
			if args[2] == failOn {
				return nil, errors.New("boom")
			}
			*deleted = append(*deleted, args[2])
			return nil, nil
		}
		t.Fatalf("unexpected op call: %v", args)
		return nil, nil
	}
	return deleted
}

func TestDeprovisioningAllSuspended(t *testing.T) {
	deleted := fakeOp(t, "")

	if err := NewOnePassword("").Deprovisioning("", false); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(*deleted, ","); got != "1,3" {
		t.Errorf("deleted %q, want 1,3", got)
	}
}

func TestDeprovisioningSingleUser(t *testing.T) {
	cases := []struct {
		email       string
		wantDeleted string
		wantErr     bool
	}{
		{"alice@example.com", "1", false},
		{"ALICE@example.com", "1", false},
		{"bob@example.com", "", false}, // not suspended: untouched
		{"nobody@example.com", "", true},
	}

	for _, c := range cases {
		deleted := fakeOp(t, "")
		err := NewOnePassword("").Deprovisioning(c.email, false)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: error = %v, wantErr %v", c.email, err, c.wantErr)
		}
		if got := strings.Join(*deleted, ","); got != c.wantDeleted {
			t.Errorf("%s: deleted %q, want %q", c.email, got, c.wantDeleted)
		}
	}
}

func TestDeprovisioningKeepsGoingAfterFailure(t *testing.T) {
	deleted := fakeOp(t, "1")

	err := NewOnePassword("").Deprovisioning("", false)
	if err == nil || !strings.Contains(err.Error(), "alice@example.com") {
		t.Errorf("expected an error mentioning alice, got %v", err)
	}
	if got := strings.Join(*deleted, ","); got != "3" {
		t.Errorf("deleted %q, want 3", got)
	}
}

func TestDeprovisioningDryRun(t *testing.T) {
	deleted := fakeOp(t, "")

	if err := NewOnePassword("").Deprovisioning("", true); err != nil {
		t.Fatal(err)
	}
	if len(*deleted) != 0 {
		t.Errorf("dry run deleted %v", *deleted)
	}
}
