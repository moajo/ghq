package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandGet(t *testing.T) {
	app := newApp()

	testCases := []struct {
		name     string
		scenario func(*testing.T, string, *_cloneArgs, *_updateArgs)
	}{
		{
			name: "simple",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

				app.Run([]string{"", "get", "motemen/ghq-test-repo"})

				expect := "https://github.com/motemen/ghq-test-repo"
				if cloneArgs.remote.String() != expect {
					t.Errorf("got: %s, expect: %s", cloneArgs.remote, expect)
				}
				if filepath.ToSlash(cloneArgs.local) != filepath.ToSlash(localDir) {
					t.Errorf("got: %s, expect: %s", filepath.ToSlash(cloneArgs.local), filepath.ToSlash(localDir))
				}
				if cloneArgs.shallow {
					t.Errorf("cloneArgs.shallow should be false")
				}
				if cloneArgs.branch != "" {
					t.Errorf("cloneArgs.branch should be empty")
				}
			},
		},
		{
			name: "-p option",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

				app.Run([]string{"", "get", "-p", "motemen/ghq-test-repo"})

				expect := "ssh://git@github.com/motemen/ghq-test-repo"
				if cloneArgs.remote.String() != expect {
					t.Errorf("got: %s, expect: %s", cloneArgs.remote, expect)
				}
				if filepath.ToSlash(cloneArgs.local) != filepath.ToSlash(localDir) {
					t.Errorf("got: %s, expect: %s", filepath.ToSlash(cloneArgs.local), filepath.ToSlash(localDir))
				}
				if cloneArgs.shallow {
					t.Errorf("cloneArgs.shallow should be false")
				}
			},
		},
		{
			name: "already cloned with -u",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")
				// mark as "already cloned", the condition may change later
				os.MkdirAll(filepath.Join(localDir, ".git"), 0755)

				app.Run([]string{"", "get", "-u", "motemen/ghq-test-repo"})

				if updateArgs.local != localDir {
					t.Errorf("got: %s, expect: %s", updateArgs.local, localDir)
				}
			},
		},
		{
			name: "shallow",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

				app.Run([]string{"", "get", "-shallow", "motemen/ghq-test-repo"})

				expect := "https://github.com/motemen/ghq-test-repo"
				if cloneArgs.remote.String() != expect {
					t.Errorf("got: %s, expect: %s", cloneArgs.remote, expect)
				}
				if filepath.ToSlash(cloneArgs.local) != filepath.ToSlash(localDir) {
					t.Errorf("got: %s, expect: %s", filepath.ToSlash(cloneArgs.local), filepath.ToSlash(localDir))
				}
				if !cloneArgs.shallow {
					t.Errorf("cloneArgs.shallow should be true")
				}
			},
		},
		{
			name: "dot slash ./",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen")
				os.MkdirAll(localDir, 0755)
				wd, _ := os.Getwd()
				os.Chdir(localDir)
				defer os.Chdir(wd)

				app.Run([]string{"", "get", "-u", "." + string(filepath.Separator) + "ghq-test-repo"})

				expect := "https://github.com/motemen/ghq-test-repo"
				if cloneArgs.remote.String() != expect {
					t.Errorf("got: %s, expect: %s", cloneArgs.remote, expect)
				}
				expectDir := filepath.Join(localDir, "ghq-test-repo")
				if cloneArgs.local != expectDir {
					t.Errorf("got: %s, expect: %s", cloneArgs.local, expectDir)
				}
			},
		},
		{
			name: "dot dot slash ../",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")
				os.MkdirAll(localDir, 0755)
				wd, _ := os.Getwd()
				os.Chdir(localDir)
				defer os.Chdir(wd)

				app.Run([]string{"", "get", "-u", ".." + string(filepath.Separator) + "ghq-another-test-repo"})

				expect := "https://github.com/motemen/ghq-another-test-repo"
				if cloneArgs.remote.String() != expect {
					t.Errorf("got: %s, expect: %s", cloneArgs.remote, expect)
				}
				expectDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-another-test-repo")
				if cloneArgs.local != expectDir {
					t.Errorf("got: %s, expect: %s", cloneArgs.local, expectDir)
				}
			},
		},
		{
			name: "specific branch",
			scenario: func(t *testing.T, tmpRoot string, cloneArgs *_cloneArgs, updateArgs *_updateArgs) {
				localDir := filepath.Join(tmpRoot, "github.com", "motemen", "ghq-test-repo")

				expectBranch := "hello"
				app.Run([]string{"", "get", "-shallow", "-b", expectBranch, "motemen/ghq-test-repo"})

				expect := "https://github.com/motemen/ghq-test-repo"
				if cloneArgs.remote.String() != expect {
					t.Errorf("got: %s, expect: %s", cloneArgs.remote, expect)
				}
				if filepath.ToSlash(cloneArgs.local) != filepath.ToSlash(localDir) {
					t.Errorf("got: %s, expect: %s", filepath.ToSlash(cloneArgs.local), filepath.ToSlash(localDir))
				}
				if cloneArgs.branch != expectBranch {
					t.Errorf("got: %q, expect: %q", cloneArgs.branch, expectBranch)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withFakeGitBackend(t, tc.scenario)
		})
	}
}
