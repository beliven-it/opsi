package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestScheduleSlots(t *testing.T) {
	slots := scheduleSlots()

	// 5 working days, 9 hours, 4 slots an hour
	if len(slots) != 5*9*4 {
		t.Errorf("got %d slots, want %d", len(slots), 5*9*4)
	}
	for _, slot := range slots {
		if slot.day == time.Saturday || slot.day == time.Sunday {
			t.Errorf("slot on a weekend day: %v", slot)
		}
		if slot.hour < scheduleFirstHour || slot.hour >= scheduleLastHour {
			t.Errorf("slot out of the time window: %v", slot)
		}
	}

	if got := (scheduleSlot{time.Wednesday, 10, 45}).cron(); got != "45 10 * * 3" {
		t.Errorf("cron = %q", got)
	}
}

func TestSlotsOfSchedule(t *testing.T) {
	cases := []struct {
		name     string
		cron     string
		timezone string
		want     []scheduleSlot
	}{
		{"weekly slot", "45 10 * * 3", "Europe/Rome", []scheduleSlot{{time.Wednesday, 10, 45}}},
		{"every working day", "0 9 * * 1-5", "Europe/Rome", []scheduleSlot{
			{time.Monday, 9, 0}, {time.Tuesday, 9, 0}, {time.Wednesday, 9, 0}, {time.Thursday, 9, 0}, {time.Friday, 9, 0}}},
		{"list of days", "30 14 * * 1,4", "Europe/Rome", []scheduleSlot{{time.Monday, 14, 30}, {time.Thursday, 14, 30}}},
		// January: Rome is UTC+1
		{"other time zone is converted", "0 8 * * 2", "Etc/UTC", []scheduleSlot{{time.Tuesday, 9, 0}}},
		{"conversion can change the day", "30 23 * * 2", "Etc/UTC", []scheduleSlot{{time.Wednesday, 0, 30}}},
		{"unknown time zone is read as UTC", "0 8 * * 2", "Nowhere/Land", []scheduleSlot{{time.Tuesday, 9, 0}}},
		{"steps cannot be placed", "*/15 * * * 1", "Europe/Rome", nil},
		{"monthly cannot be placed", "0 9 1 * *", "Europe/Rome", nil},
		{"garbage", "not a cron", "Europe/Rome", nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := slotsOfSchedule(c.cron, c.timezone)
			if len(got) != len(c.want) {
				t.Fatalf("got %v, want %v", got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("got %v, want %v", got, c.want)
				}
			}
		})
	}
}

func TestPickAvoidsLoadedSlots(t *testing.T) {
	load := scheduleLoad{}

	// Fill every slot once, but leave one empty
	free := scheduleSlot{time.Thursday, 15, 30}
	for _, slot := range scheduleSlots() {
		if slot != free {
			load[slot] = 1
		}
	}

	if got := load.pick("any/project"); got != free {
		t.Errorf("picked %v, want the only free slot %v", got, free)
	}

	// Once the free slot is taken too, every slot has the same load: the pick
	// is anywhere, but always the same for the same project
	load.add(free.cron(), scheduleTimezone)
	first := load.pick("a/project")
	if again := load.pick("a/project"); again != first {
		t.Errorf("same project got %v then %v", first, again)
	}
}

func TestPickSpreadsProjects(t *testing.T) {
	load := scheduleLoad{}
	used := map[scheduleSlot]int{}

	// Fewer projects than slots: nobody should share a slot
	for i := 0; i < 100; i++ {
		slot := load.pick(fmt.Sprintf("group/project-%d", i))
		used[slot]++
		load.add(slot.cron(), scheduleTimezone)
	}

	for slot, n := range used {
		if n > 1 {
			t.Errorf("slot %v used %d times with free slots left", slot, n)
		}
	}
}

func TestSlotKeyRoundsDown(t *testing.T) {
	if got := slotKey(scheduleSlot{time.Monday, 10, 20}); got != (scheduleSlot{time.Monday, 10, 15}) {
		t.Errorf("slotKey = %v", got)
	}
}

type scheduleServer struct {
	mu      sync.Mutex
	created []map[string]any
	t       *testing.T
}

func (s *scheduleServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/projects/corporate/with-schedule" || r.URL.EscapedPath() == "/projects/corporate%2Fwith-schedule":
			fmt.Fprint(w, `{"id":2,"path_with_namespace":"corporate/with-schedule","default_branch":"main"}`)
		case r.URL.EscapedPath() == "/projects/corporate%2Fstaging-only":
			fmt.Fprint(w, `{"id":3,"path_with_namespace":"corporate/staging-only","default_branch":"staging"}`)
		case r.URL.EscapedPath() == "/projects/corporate%2Fnew":
			fmt.Fprint(w, `{"id":1,"path_with_namespace":"corporate/new","default_branch":"main"}`)
		case r.URL.Path == "/projects" && r.Method == "GET":
			fmt.Fprint(w, `[{"id":1},{"id":2},{"id":3}]`)
		case r.URL.Path == "/projects/2/pipeline_schedules" && r.Method == "GET":
			fmt.Fprint(w, `[{"id":9,"description":"x","ref":"refs/heads/main","cron":"0 9 * * 1-5","cron_timezone":"Europe/Rome","active":true}]`)
		case strings.HasSuffix(r.URL.Path, "/pipeline_schedules") && r.Method == "GET":
			fmt.Fprint(w, `[]`)
		case r.URL.Path == "/projects/3/repository/branches/main":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message":"404 Branch Not Found"}`)
		case strings.Contains(r.URL.Path, "/repository/branches/main"):
			fmt.Fprint(w, `{"name":"main"}`)
		case strings.HasSuffix(r.URL.Path, "/pipeline_schedules") && r.Method == "POST":
			var body map[string]any
			json.NewDecoder(r.Body).Decode(&body)
			body["project"] = r.URL.Path
			s.mu.Lock()
			s.created = append(s.created, body)
			s.mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id":1}`)
		case strings.HasPrefix(r.URL.EscapedPath(), "/projects/missing"):
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"message":"404 Project Not Found"}`)
		default:
			s.t.Errorf("unexpected request %s %s", r.Method, r.URL.EscapedPath())
		}
	})
}

func TestCreateSchedule(t *testing.T) {
	s := &scheduleServer{t: t}
	server := httptest.NewServer(s.handler())
	defer server.Close()

	g := &gitlab{apiURL: server.URL, token: "t"}

	err := g.CreateSchedule([]string{"corporate/new", "corporate/with-schedule", "corporate/staging-only", "missing/project"})
	if err == nil || !strings.Contains(err.Error(), "missing/project") || strings.Contains(err.Error(), "corporate/new") {
		t.Errorf("only the missing project should fail, got %v", err)
	}

	if len(s.created) != 2 {
		t.Fatalf("created %d schedules, want 2 (the project with an active schedule is skipped): %v", len(s.created), s.created)
	}

	byProject := map[string]map[string]any{}
	for _, c := range s.created {
		byProject[c["project"].(string)] = c
	}

	for project, wantRef := range map[string]string{
		"/projects/1/pipeline_schedules": "main",
		"/projects/3/pipeline_schedules": "staging", // no main branch: default branch
	} {
		c := byProject[project]
		if c == nil {
			t.Fatalf("no schedule created on %s", project)
		}
		if c["ref"] != wantRef {
			t.Errorf("%s: ref = %v, want %v", project, c["ref"], wantRef)
		}
		if c["description"] != scheduleDescription || c["cron_timezone"] != scheduleTimezone || c["active"] != true {
			t.Errorf("%s: unexpected payload %v", project, c)
		}
	}

	// The two new schedules must not share a slot
	if byProject["/projects/1/pipeline_schedules"]["cron"] == byProject["/projects/3/pipeline_schedules"]["cron"] {
		t.Errorf("the two schedules got the same cron %v", byProject["/projects/1/pipeline_schedules"]["cron"])
	}
}

func TestUpdateSchedule(t *testing.T) {
	var mu sync.Mutex
	taken := []string{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/user":
			fmt.Fprint(w, `{"id":10,"username":"cicd","state":"active"}`)
		case r.URL.EscapedPath() == "/projects/corporate%2Fapp":
			fmt.Fprint(w, `{"id":1,"path_with_namespace":"corporate/app","default_branch":"main"}`)
		case r.URL.EscapedPath() == "/projects/corporate%2Fempty":
			fmt.Fprint(w, `{"id":2,"path_with_namespace":"corporate/empty","default_branch":"main"}`)
		case r.URL.Path == "/projects/1/pipeline_schedules":
			fmt.Fprint(w, `[
			  {"id":1,"cron":"0 9 * * 1","ref":"refs/heads/main","active":true,"owner":{"id":10,"username":"cicd","state":"active"}},
			  {"id":2,"cron":"0 10 * * 2","ref":"refs/heads/main","active":true,"owner":{"id":20,"username":"gone","state":"blocked"}},
			  {"id":3,"cron":"0 11 * * 3","ref":"refs/heads/main","active":false,"owner":null},
			  {"id":4,"cron":"0 12 * * 4","ref":"refs/heads/main","active":true,"owner":{"id":30,"username":"other","state":"active"}}]`)
		case r.URL.Path == "/projects/2/pipeline_schedules":
			fmt.Fprint(w, `[]`)
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/take_ownership"):
			if r.URL.Path == "/projects/1/pipeline_schedules/4/take_ownership" {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, `{"message":"403 Forbidden"}`)
				return
			}
			mu.Lock()
			taken = append(taken, r.URL.Path)
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.EscapedPath())
		}
	}))
	defer server.Close()

	g := &gitlab{apiURL: server.URL, token: "t"}

	// A dry run must not change anything
	if err := g.UpdateSchedule([]string{"corporate/app", "corporate/empty"}, true); err != nil {
		t.Errorf("dry run failed: %v", err)
	}
	if len(taken) != 0 {
		t.Errorf("dry run took ownership of %v", taken)
	}

	// Schedules already owned are left alone, the others are taken
	// (including the one with no owner and the inactive one); a failure
	// is reported without stopping the rest
	err := g.UpdateSchedule([]string{"corporate/app", "corporate/empty"}, false)
	if err == nil || !strings.Contains(err.Error(), "corporate/app") {
		t.Errorf("expected the failure on schedule 4, got %v", err)
	}

	want := "/projects/1/pipeline_schedules/2/take_ownership,/projects/1/pipeline_schedules/3/take_ownership"
	if got := strings.Join(taken, ","); got != want {
		t.Errorf("took ownership of %q, want %q", got, want)
	}
}

func TestDescribeSchedule(t *testing.T) {
	got := describeSchedule(gitlabPipelineSchedule{Cron: "0 9 * * 1", Ref: "refs/heads/main", Active: true})
	if got != "0 9 * * 1 on main, active, owner none" {
		t.Errorf("got %q", got)
	}
}
