package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	content := "A=1\n# comment\nB = \"two\"\nC='3'\nD\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("A", "real") // existing env must win

	LoadDotEnv(path)

	if got := os.Getenv("A"); got != "real" {
		t.Fatalf("A = %q, want existing env value to win", got)
	}
	for k, want := range map[string]string{"B": "two", "C": "3"} {
		if got := os.Getenv(k); got != want {
			t.Fatalf("%s = %q, want %q", k, got, want)
		}
	}
	if got := os.Getenv("D"); got != "" {
		t.Fatalf("D = %q, want empty (no value)", got)
	}
}

func TestLoadDotEnvMissingFile(t *testing.T) {
	LoadDotEnv(filepath.Join(t.TempDir(), "nope")) // must not panic
}
