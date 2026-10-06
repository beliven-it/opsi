package gitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	// The time zones are needed to read the existing schedules, also on machines without tzdata
	_ "time/tzdata"
)

// The default pipeline schedule is weekly: one working day and one time of the
// day for each project. The slot is chosen among the ones with the fewest
// schedules, so the pipelines are spread instead of starting all together.
const (
	scheduleDescription = "opsi default schedule"
	scheduleTimezone    = "Europe/Rome"
	scheduleFirstHour   = 9  // first slot at 09:00
	scheduleLastHour    = 18 // no slot starting at or after 18:00
	scheduleSlotMinutes = 15
	scheduleBranch      = "main"
)

var scheduleDays = []time.Weekday{
	time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
}

// scheduleSlot is a weekly time, in the scheduleTimezone.
type scheduleSlot struct {
	day    time.Weekday
	hour   int
	minute int
}

func (s scheduleSlot) cron() string {
	return fmt.Sprintf("%d %d * * %d", s.minute, s.hour, int(s.day))
}

func (s scheduleSlot) String() string {
	return fmt.Sprintf("%s %02d:%02d", strings.ToLower(s.day.String()), s.hour, s.minute)
}

// scheduleSlots returns every slot a schedule can be placed on.
func scheduleSlots() []scheduleSlot {
	slots := []scheduleSlot{}
	for _, day := range scheduleDays {
		for hour := scheduleFirstHour; hour < scheduleLastHour; hour++ {
			for minute := 0; minute < 60; minute += scheduleSlotMinutes {
				slots = append(slots, scheduleSlot{day, hour, minute})
			}
		}
	}

	return slots
}

// cronField expands a cron field made of numbers, lists, ranges and "*"
// into the values it matches, within [min, max]. It returns false for
// anything fancier (steps, names), which cannot be placed on a slot.
func cronField(field string, min int, max int) ([]int, bool) {
	values := []int{}

	for _, part := range strings.Split(field, ",") {
		switch {
		case part == "*":
			for v := min; v <= max; v++ {
				values = append(values, v)
			}
		case strings.Contains(part, "-"):
			bounds := strings.SplitN(part, "-", 2)
			from, err1 := strconv.Atoi(bounds[0])
			to, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || from < min || to > max || from > to {
				return nil, false
			}
			for v := from; v <= to; v++ {
				values = append(values, v)
			}
		default:
			v, err := strconv.Atoi(part)
			if err != nil || v < min || v > max {
				return nil, false
			}
			values = append(values, v)
		}
	}

	return values, true
}

// slotsOfSchedule returns the slots, in the scheduleTimezone, an existing
// schedule runs on. Schedules that run more than once a day are counted on
// each of the matching slots; the ones that cannot be understood are ignored.
func slotsOfSchedule(cron string, timezone string) []scheduleSlot {
	fields := strings.Fields(cron)
	if len(fields) != 5 || fields[2] != "*" || fields[3] != "*" {
		return nil
	}

	minutes, ok1 := cronField(fields[0], 0, 59)
	hours, ok2 := cronField(fields[1], 0, 23)
	days, ok3 := cronField(strings.ReplaceAll(fields[4], "7", "0"), 0, 6)
	if !ok1 || !ok2 || !ok3 {
		return nil
	}

	source, err := time.LoadLocation(timezone)
	if err != nil {
		source = time.UTC
	}
	target, _ := time.LoadLocation(scheduleTimezone)
	if target == nil {
		target = time.UTC
	}

	// A fixed reference week (starting on a Sunday) to convert the time zone
	referenceSunday := time.Date(2026, time.January, 4, 0, 0, 0, 0, source)

	slots := []scheduleSlot{}
	for _, day := range days {
		for _, hour := range hours {
			for _, minute := range minutes {
				at := referenceSunday.AddDate(0, 0, day).Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute).In(target)
				slots = append(slots, scheduleSlot{at.Weekday(), at.Hour(), at.Minute()})
			}
		}
	}

	return slots
}

// scheduleLoad counts how many schedules run on each slot.
type scheduleLoad map[scheduleSlot]int

// slotKey is the slot a time belongs to: the minute is rounded down to the
// slot size, so a schedule at 10:20 weighs on the 10:15 slot.
func slotKey(s scheduleSlot) scheduleSlot {
	s.minute = s.minute - s.minute%scheduleSlotMinutes
	return s
}

func (l scheduleLoad) add(cron string, timezone string) {
	for _, slot := range slotsOfSchedule(cron, timezone) {
		l[slotKey(slot)]++
	}
}

// pick returns the least loaded slot. Among the equally loaded ones it picks
// by the hash of the seed (the project path), so the choice does not always
// fall on the first slot of the week but stays the same for the same project.
func (l scheduleLoad) pick(seed string) scheduleSlot {
	slots := scheduleSlots()

	lowest := -1
	candidates := []scheduleSlot{}
	for _, slot := range slots {
		load := l[slot]
		switch {
		case lowest == -1 || load < lowest:
			lowest = load
			candidates = []scheduleSlot{slot}
		case load == lowest:
			candidates = append(candidates, slot)
		}
	}

	hash := fnv.New32a()
	hash.Write([]byte(seed))

	return candidates[int(hash.Sum32())%len(candidates)]
}

type gitlabScheduleOwner struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	State    string `json:"state"`
}

type gitlabPipelineSchedule struct {
	ID           int                  `json:"id"`
	Description  string               `json:"description"`
	Ref          string               `json:"ref"`
	Cron         string               `json:"cron"`
	CronTimezone string               `json:"cron_timezone"`
	Active       bool                 `json:"active"`
	Owner        *gitlabScheduleOwner `json:"owner"`
	LastPipeline *struct {
		Status string `json:"status"`
	} `json:"last_pipeline"`
}

// listSchedules returns the pipeline schedules of a project.
func (g *gitlab) listSchedules(projectID int) ([]gitlabPipelineSchedule, error) {
	schedules := []gitlabPipelineSchedule{}

	for page := 1; ; page++ {
		response, err := g.request("GET", fmt.Sprintf("/projects/%d/pipeline_schedules", projectID), nil, map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": "100",
		})
		if err != nil {
			return nil, err
		}

		var list []gitlabPipelineSchedule
		if err := json.Unmarshal(response, &list); err != nil {
			return nil, err
		}

		schedules = append(schedules, list...)
		if len(list) < 100 {
			return schedules, nil
		}
	}
}

// loadSchedules counts the active schedules of every project. The projects
// that cannot be read are ignored: the load is only used to choose a slot.
func (g *gitlab) loadSchedules() (scheduleLoad, error) {
	projects, err := g.listProjects()
	if err != nil {
		return nil, err
	}

	load := scheduleLoad{}

	var (
		mu   sync.Mutex
		wg   sync.WaitGroup
		done int
	)
	slots := make(chan struct{}, deprovisioningConcurrency)

	for _, project := range projects {
		wg.Add(1)
		slots <- struct{}{}
		go func(id int) {
			defer wg.Done()
			defer func() { <-slots }()

			schedules, _ := g.listSchedules(id)

			mu.Lock()
			defer mu.Unlock()
			for _, schedule := range schedules {
				if schedule.Active {
					load.add(schedule.Cron, schedule.CronTimezone)
				}
			}
			done++
			progress("Checking existing schedules... %d/%d", done, len(projects))
		}(project.ID)
	}

	wg.Wait()
	endProgress()

	return load, nil
}

// scheduleRef is the branch the schedule runs on: the production branch
// when the project has it, otherwise the default branch of the project.
func (g *gitlab) scheduleRef(projectID int, defaultBranch string) string {
	_, err := g.request("GET", fmt.Sprintf("/projects/%d/repository/branches/%s", projectID, scheduleBranch), nil, nil)
	if err == nil {
		return scheduleBranch
	}

	return defaultBranch
}

// createSchedule creates the default schedule of a project on the least
// loaded slot and counts it in the load. It returns a description of it.
func (g *gitlab) createSchedule(projectID int, path string, defaultBranch string, load scheduleLoad) (string, error) {
	ref := g.scheduleRef(projectID, defaultBranch)
	slot := load.pick(path)

	_, err := g.request("POST", fmt.Sprintf("/projects/%d/pipeline_schedules", projectID), map[string]any{
		"description":   scheduleDescription,
		"ref":           ref,
		"cron":          slot.cron(),
		"cron_timezone": scheduleTimezone,
		"active":        true,
	}, nil)
	if err != nil {
		return "", err
	}

	load.add(slot.cron(), scheduleTimezone)

	return fmt.Sprintf("%s, on %s", slot, ref), nil
}

type scheduleProject struct {
	ID                int    `json:"id"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
}

func (g *gitlab) findProject(ref string) (scheduleProject, error) {
	var project scheduleProject

	response, err := g.request("GET", "/projects/"+url.PathEscape(ref), nil, nil)
	if err != nil {
		return project, err
	}

	err = json.Unmarshal(response, &project)

	return project, err
}

// CreateSchedule creates the default pipeline schedule on each of the
// projects (numeric ID or full path). A project that already has an active
// schedule is left alone.
func (g *gitlab) CreateSchedule(projects []string) error {
	var failures []error
	var load scheduleLoad

	for _, ref := range projects {
		project, err := g.findProject(ref)
		if err != nil {
			fmt.Printf("Error when reading %s: %v\n", ref, err)
			failures = append(failures, fmt.Errorf("%s: %w", ref, err))
			continue
		}

		existing, err := g.listSchedules(project.ID)
		if err != nil {
			fmt.Printf("Error when reading the schedules of %s: %v\n", project.PathWithNamespace, err)
			failures = append(failures, fmt.Errorf("%s: %w", project.PathWithNamespace, err))
			continue
		}

		var active *gitlabPipelineSchedule
		for i := range existing {
			if existing[i].Active {
				active = &existing[i]
				break
			}
		}
		if active != nil {
			fmt.Printf("Skipped %s: it already has an active schedule (%s, on %s)\n", project.PathWithNamespace, active.Cron, active.Ref)
			continue
		}

		// The load is read once, and only when there is something to create
		if load == nil {
			load, err = g.loadSchedules()
			if err != nil {
				return err
			}
		}

		description, err := g.createSchedule(project.ID, project.PathWithNamespace, project.DefaultBranch, load)
		if err != nil {
			fmt.Printf("Error when creating the schedule for %s: %v\n", project.PathWithNamespace, err)
			failures = append(failures, fmt.Errorf("%s: %w", project.PathWithNamespace, err))
			continue
		}

		fmt.Printf("Schedule created for %s: %s\n", project.PathWithNamespace, description)
	}

	return errors.Join(failures...)
}
