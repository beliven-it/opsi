package gitlab

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMirrorProjectName(t *testing.T) {
	cases := []struct {
		url     string
		want    string
		wantErr bool
	}{
		{"https://user:token@gitlab.example.com/group/my-repo.git", "my-repo", false},
		{"https://gitlab.example.com/my-repo.git", "my-repo", false},
		{"https://gitlab.example.com/group/no-suffix", "", true},
		{"", "", true},
	}

	for _, c := range cases {
		got, err := mirrorProjectName(c.url)
		if (err != nil) != c.wantErr {
			t.Errorf("mirrorProjectName(%q) error = %v, wantErr %v", c.url, err, c.wantErr)
		}
		if got != c.want {
			t.Errorf("mirrorProjectName(%q) = %q, want %q", c.url, got, c.want)
		}
	}
}

func TestEnableMirrorWithRetry(t *testing.T) {
	mirrorRetryDelay = time.Millisecond
	defer func() { mirrorRetryDelay = 2 * time.Second }()

	cases := []struct {
		name         string
		failures     int
		wantErr      bool
		wantAttempts int
	}{
		{"succeeds at first attempt", 0, false, 1},
		{"recovers after transient failures", 2, false, 3},
		{"gives up after max attempts", 10, true, mirrorMaxAttempts},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var mu sync.Mutex
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				calls++
				n := calls
				mu.Unlock()
				if n <= c.failures {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, `{"message":"boom"}`)
					return
				}
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `{}`)
			}))
			defer server.Close()

			g := &gitlab{apiURL: server.URL, token: "t"}
			g.mirror.GroupPath = "mirror.example.com/group"

			err := g.enableMirrorWithRetry(1, "repo")
			if (err != nil) != c.wantErr {
				t.Errorf("error = %v, wantErr %v", err, c.wantErr)
			}
			if calls != c.wantAttempts {
				t.Errorf("attempts = %d, want %d", calls, c.wantAttempts)
			}
		})
	}
}

// A mirror with a malformed URL must be skipped without being deleted.
func TestUpdateMirroringSkipsMalformedURL(t *testing.T) {
	var mu sync.Mutex
	var deletes []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/projects/1/remote_mirrors"):
			fmt.Fprint(w, `[{"id":9,"enabled":true,"url":"https://host/not-a-git-url"}]`)
		case r.Method == "GET" && r.URL.Path == "/projects":
			fmt.Fprint(w, `[{"id":1}]`)
		case r.Method == "DELETE":
			deletes = append(deletes, r.URL.Path)
		default:
			fmt.Fprint(w, `[]`)
		}
	}))
	defer server.Close()

	g := &gitlab{apiURL: server.URL, token: "t"}
	err := g.UpdateMirroring()
	if err == nil {
		t.Error("expected an error for the malformed mirror url")
	}
	if len(deletes) != 0 {
		t.Errorf("mirror must not be deleted, got DELETE on %v", deletes)
	}
}

func TestCheckMirroringExistence(t *testing.T) {
	cases := []struct {
		name    string
		body    string
		wantID  int
		wantHas bool
	}{
		{"no mirror", `[]`, 0, false},
		{"enabled mirror", `[{"id":1,"enabled":true,"url":"u"}]`, 1, true},
		{"disabled mirror is still a mirror", `[{"id":2,"enabled":false,"url":"u"}]`, 2, true},
		{"enabled is preferred over disabled", `[{"id":3,"enabled":false},{"id":4,"enabled":true}]`, 4, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, c.body)
			}))
			defer server.Close()

			g := &gitlab{apiURL: server.URL, token: "t"}
			mirror, has, err := g.checkMirroringExistence(1)
			if err != nil {
				t.Fatal(err)
			}
			if has != c.wantHas || mirror.ID != c.wantID {
				t.Errorf("got (%d, %v), want (%d, %v)", mirror.ID, has, c.wantID, c.wantHas)
			}
		})
	}
}

// newMirrorFixture serves a source instance with one project and a mirror
// instance, recording the calls that matter.
func newMirrorFixture(t *testing.T, sourceMirrors string, destinationExists bool) (*gitlab, *[]string) {
	t.Helper()

	var mu sync.Mutex
	calls := []string{}
	record := func(r *http.Request) {
		mu.Lock()
		calls = append(calls, r.Method+" "+r.URL.EscapedPath())
		mu.Unlock()
	}

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record(r)
		switch {
		case r.Method == "GET" && r.URL.Path == "/projects/7":
			fmt.Fprint(w, `{"id":7,"name":"Wiki","path":"wiki","path_with_namespace":"corporate/wiki"}`)
		case r.Method == "GET" && r.URL.Path == "/projects/7/remote_mirrors":
			fmt.Fprint(w, sourceMirrors)
		case r.Method == "POST" && r.URL.Path == "/projects/7/remote_mirrors":
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message":"404"}`)
		}
	}))
	t.Cleanup(source.Close)

	mirror := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record(r)
		switch {
		case r.Method == "GET" && r.URL.EscapedPath() == "/projects/mirrors%2Fwiki" && destinationExists:
			fmt.Fprint(w, `{"id":99}`)
		case r.Method == "POST" && r.URL.Path == "/projects":
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id":99}`)
		case r.Method == "PATCH":
			fmt.Fprint(w, `{}`)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message":"404"}`)
		}
	}))
	t.Cleanup(mirror.Close)

	g := &gitlab{apiURL: source.URL, token: "t"}
	g.mirror.ApiURL = mirror.URL
	g.mirror.GroupPath = "gitlab.example.com/mirrors"
	g.mirror.GroupID = 5

	return g, &calls
}

func TestCreateMirrorWithExistingDestination(t *testing.T) {
	g, calls := newMirrorFixture(t, `[]`, true)

	if err := g.CreateMirror([]int{7}); err != nil {
		t.Fatal(err)
	}

	joined := strings.Join(*calls, "\n")
	if strings.Contains(joined, "POST /projects\n") || strings.HasSuffix(joined, "POST /projects") {
		t.Errorf("destination exists, it must not be created again:\n%s", joined)
	}
	if !strings.Contains(joined, "POST /projects/7/remote_mirrors") {
		t.Errorf("remote mirror was not created:\n%s", joined)
	}
}

func TestCreateMirrorCreatesMissingDestination(t *testing.T) {
	g, calls := newMirrorFixture(t, `[]`, false)

	if err := g.CreateMirror([]int{7}); err != nil {
		t.Fatal(err)
	}

	joined := strings.Join(*calls, "\n")
	if !strings.Contains(joined, "POST /projects\n") && !strings.HasSuffix(joined, "POST /projects") {
		t.Errorf("missing destination was not created:\n%s", joined)
	}
	if !strings.Contains(joined, "POST /projects/7/remote_mirrors") {
		t.Errorf("remote mirror was not created:\n%s", joined)
	}
}

func TestCreateMirrorRefusesProjectsWithAMirror(t *testing.T) {
	g, calls := newMirrorFixture(t, `[{"id":1,"enabled":false,"url":"u"}]`, true)

	if err := g.CreateMirror([]int{7}); err == nil {
		t.Error("expected an error for a project that already has a mirror")
	}
	if strings.Contains(strings.Join(*calls, "\n"), "POST /projects/7/remote_mirrors") {
		t.Error("no mirror must be created for a project that has one")
	}
}

func TestCreateMirrorKeepsGoingAfterAFailure(t *testing.T) {
	g, calls := newMirrorFixture(t, `[]`, true)

	err := g.CreateMirror([]int{404, 7})
	if err == nil {
		t.Error("expected the failure to be reported")
	}
	if !strings.Contains(strings.Join(*calls, "\n"), "POST /projects/7/remote_mirrors") {
		t.Error("the second project must still be processed")
	}
}
