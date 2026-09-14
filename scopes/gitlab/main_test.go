package gitlab

import (
	"testing"
)

func TestSameValue(t *testing.T) {
	cases := []struct {
		name     string
		current  any
		desired  any
		expected bool
	}{
		// The group API answers null for a boolean that was never set
		{"unset boolean matches false", nil, false, true},
		{"unset boolean differs from true", nil, true, false},
		{"unset value never matches a string", nil, "maintainer", false},
		{"matching booleans", true, true, true},
		{"differing booleans", false, true, false},
		// Both sides of the diff come from JSON, so numbers are float64
		{"matching numbers", float64(48), float64(48), true},
		{"differing numbers", float64(24), float64(48), false},
		{"matching strings", "maintainer", "maintainer", true},
		{"differing strings", "developer", "maintainer", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := sameValue(c.current, c.desired); got != c.expected {
				t.Errorf("sameValue(%v, %v) = %v, want %v", c.current, c.desired, got, c.expected)
			}
		})
	}
}

func TestCurrentValue(t *testing.T) {
	if got := currentValue(nil, true); got != false {
		t.Errorf("an unset boolean should read as false, got %v", got)
	}

	if got := currentValue(nil, "maintainer"); got != nil {
		t.Errorf("an unset string should stay nil, got %v", got)
	}

	if got := currentValue("developer", "maintainer"); got != "developer" {
		t.Errorf("a set value should be left alone, got %v", got)
	}
}

func TestDesiredGroupSettings(t *testing.T) {
	desired, err := desiredGroupSettings()
	if err != nil {
		t.Fatal(err)
	}

	if len(desired) != 6 {
		t.Errorf("expected 6 managed settings, got %d: %v", len(desired), desired)
	}

	// The JSON round trip is what keeps the diff comparable with the API answer
	if grace, ok := desired["two_factor_grace_period"].(float64); !ok || grace != 48 {
		t.Errorf("two_factor_grace_period should be float64(48), got %#v", desired["two_factor_grace_period"])
	}

	if desired["mentions_disabled"] != true {
		t.Errorf("mentions_disabled should be true, got %#v", desired["mentions_disabled"])
	}

	if desired["project_creation_level"] != "maintainer" {
		t.Errorf("project_creation_level should be maintainer, got %#v", desired["project_creation_level"])
	}
}

func TestGroupSettingsDiff(t *testing.T) {
	desired, err := desiredGroupSettings()
	if err != nil {
		t.Fatal(err)
	}

	conformant := map[string]any{"id": float64(1), "full_path": "conformant"}
	for key, value := range desired {
		conformant[key] = value
	}

	if changes, _ := groupSettingsDiff(conformant, desired); len(changes) != 0 {
		t.Errorf("a conformant group should not drift, got %v", changes)
	}

	drifting := map[string]any{}
	for key, value := range conformant {
		drifting[key] = value
	}
	drifting["mentions_disabled"] = nil
	drifting["project_creation_level"] = "developer"

	changes, from := groupSettingsDiff(drifting, desired)
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %v", changes)
	}

	if changes["mentions_disabled"] != true || changes["project_creation_level"] != "maintainer" {
		t.Errorf("unexpected changes: %v", changes)
	}

	// The report must show what the group page shows, not the raw null
	if from["mentions_disabled"] != false {
		t.Errorf("an unset boolean should be reported as false, got %#v", from["mentions_disabled"])
	}

	if from["project_creation_level"] != "developer" {
		t.Errorf("unexpected previous value: %#v", from["project_creation_level"])
	}
}

func TestFilterGroupHierarchy(t *testing.T) {
	groups := []map[string]any{
		{"id": float64(1), "full_path": "corporate"},
		{"id": float64(2), "full_path": "corporate/dev"},
		{"id": float64(3), "full_path": "corporate/dev/packages"},
		{"id": float64(4), "full_path": "corporate-legacy"},
		{"id": float64(5), "full_path": "biofarma"},
	}

	if got := filterGroupHierarchy(groups, 0); len(got) != len(groups) {
		t.Errorf("a zero groupID should keep every group, got %d", len(got))
	}

	got := filterGroupHierarchy(groups, 1)
	if len(got) != 3 {
		t.Fatalf("expected the group and its 2 descendants, got %d: %v", len(got), got)
	}

	// A sibling sharing the prefix must not be swept in
	for _, group := range got {
		if _, fullPath := groupIdentity(group); fullPath == "corporate-legacy" {
			t.Error("a sibling sharing the path prefix should be left out")
		}
	}

	if got := filterGroupHierarchy(groups, 99); got != nil {
		t.Errorf("an unknown groupID should match nothing, got %v", got)
	}
}

func TestGroupIdentity(t *testing.T) {
	id, fullPath := groupIdentity(map[string]any{"id": float64(191), "full_path": "opsi-test"})
	if id != 191 || fullPath != "opsi-test" {
		t.Errorf("got (%d, %q)", id, fullPath)
	}

	// A malformed answer must not panic
	if id, fullPath := groupIdentity(map[string]any{}); id != 0 || fullPath != "" {
		t.Errorf("got (%d, %q)", id, fullPath)
	}
}
