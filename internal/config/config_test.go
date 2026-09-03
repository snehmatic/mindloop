package config

import (
	"path/filepath"
	"sync"
	"testing"
)

func TestWriteToYAMLErrorRoundTripIsAtomic(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(t.TempDir())

	want := UserConfig{
		Name:            "atomic-user",
		Mode:            "local",
		EditorWideWidth: true,
		PointsConfig:    PointsConfig{Task: 7},
	}
	if err := want.WriteToYAMLError(); err != nil {
		t.Fatalf("write user config: %v", err)
	}

	var got UserConfig
	if err := got.ReadFromYAML(); err != nil {
		t.Fatalf("read user config: %v", err)
	}
	if got.Name != want.Name || got.Mode != want.Mode || !got.EditorWideWidth || got.PointsConfig.Task != want.PointsConfig.Task {
		t.Fatalf("round trip mismatch: got %+v, want %+v", got, want)
	}
	if matches, err := filepath.Glob(".user_config.yaml.tmp-*"); err != nil {
		t.Fatalf("find temporary files: %v", err)
	} else if len(matches) != 0 {
		t.Fatalf("temporary files remain after atomic write: %v", matches)
	}
}

func TestUpdateUserConfigSerializesReadModifyWrite(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(t.TempDir())
	if err := (UserConfig{Mode: "local"}).WriteToYAMLError(); err != nil {
		t.Fatalf("seed user config: %v", err)
	}
	var baseline UserConfig
	if err := baseline.ReadFromYAML(); err != nil {
		t.Fatalf("read seeded user config: %v", err)
	}

	const updates = 32
	var wg sync.WaitGroup
	errs := make(chan error, updates)
	for range updates {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- UpdateUserConfig(func(uc *UserConfig) error {
				uc.PointsConfig.Task++
				return nil
			})
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("update user config: %v", err)
		}
	}

	var got UserConfig
	if err := got.ReadFromYAML(); err != nil {
		t.Fatalf("read final user config: %v", err)
	}
	wantTaskPoints := baseline.PointsConfig.Task + updates
	if got.PointsConfig.Task != wantTaskPoints {
		t.Fatalf("lost concurrent updates: got task points %d, want %d", got.PointsConfig.Task, wantTaskPoints)
	}
}
