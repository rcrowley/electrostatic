package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMainSchedule checks that, at a fixed reference time, electrostatic builds
// only the page that is currently published: not the draft, not the
// not-yet-published page, and not the expired one.
func TestMainSchedule(t *testing.T) {
	out := t.TempDir()
	Main([]string{
		"electrostatic",
		"-l", "testdata/layout.html",
		"-n", "2026-06-04 00:00:00",
		"-o", out,
		"testdata/src",
	}, os.Stdin, os.Stdout)

	if _, err := os.Stat(filepath.Join(out, "live", "index.html")); err != nil {
		t.Errorf("expected the live page to be built: %v", err)
	}
	for _, name := range []string{"draft", "scheduled", "expired"} {
		if _, err := os.Stat(filepath.Join(out, name, "index.html")); !os.IsNotExist(err) {
			t.Errorf("%s page should not have been built (err=%v)", name, err)
		}
	}
}

// TestMainScheduleArrives checks that advancing the reference time past a
// page's publish date causes it to be built.
func TestMainScheduleArrives(t *testing.T) {
	out := t.TempDir()
	Main([]string{
		"electrostatic",
		"-l", "testdata/layout.html",
		"-n", "2027-01-01 00:00:00",
		"-o", out,
		"testdata/src",
	}, os.Stdin, os.Stdout)

	if _, err := os.Stat(filepath.Join(out, "scheduled", "index.html")); err != nil {
		t.Errorf("scheduled page should be built once its time has passed: %v", err)
	}
}
